package initialize

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
	"prosto-delaj-api/models"
	"reflect"
	"strings"
	"time"

	"github.com/caarlos0/env/v6"
	vault "github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
	"github.com/joho/godotenv"
)

type customFormatter struct{}

const (
	envFilePath          = `.env`
	vaultSessionFilePath = `wrapped_vault_session_token`
	vaultPathFilePath    = `wrapped_vault_paths_token`
	configFilePath       = "configs/config.json"
)

func LoadConfiguration(appState *models.AppState) error {
	if err := loadEnvironment(appState.Env); err != nil {
		return err
	}
	if appState.Env.IsReadConfig {
		logrus.Info("load local config")
		if err := loadConfig(appState.ConfigService); err != nil {
			return err
		}
		logrus.Info("load local config success")
	} else {
		logrus.Info("load vault")
		if err := loadVault(appState); err != nil {
			return err
		}
	}
	return nil
}

func loadVault(appState *models.AppState) error {
	// Чтение данных из файлов
	wrappedSecretID, err := readVaultFile(vaultSessionFilePath)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}

	wrappedSecretPaths, err := readVaultFile(vaultPathFilePath)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}

	appState.ConfigVault = &models.ConfigVault{
		WrappedSecretID:    wrappedSecretID,
		WrappedSecretPaths: wrappedSecretPaths,
	}

	client, err := vault.New(
		vault.WithAddress(appState.Env.VaultAddr),
	)
	if err != nil {
		logrus.Error(err)
		return err
	}

	err = client.SetToken(strings.TrimSpace(appState.ConfigVault.WrappedSecretPaths))
	if err != nil {
		logrus.Error(err)
		return err
	}
	unwrapPath, err := client.System.Unwrap(context.Background(), schema.UnwrapRequest{})
	if err != nil {
		logrus.Error(err.Error())
		return err
	}

	data := unwrapPath.Data["data"]
	decodedData := strings.ReplaceAll(data.(string), "\\n", "\n")
	decodedData = strings.ReplaceAll(decodedData, "\\\"", "\"")
	var pathConfig []models.VaultSecretsConfig
	err = yaml.Unmarshal([]byte(decodedData), &pathConfig)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	err = client.SetToken(strings.TrimSpace(appState.ConfigVault.WrappedSecretID))
	if err != nil {
		logrus.Error(err)
		return err
	}
	// Распаковка WrappedSecretID
	unwrapID, err := client.System.Unwrap(context.Background(), schema.UnwrapRequest{})
	if err != nil {
		logrus.Error(err.Error())
		return err
	}

	data2 := unwrapID.Data["secret_id"]

	// Настройка AppRole авторизации
	authAppRole, err := client.Auth.AppRoleLogin(context.Background(),
		schema.AppRoleLoginRequest{
			RoleId:   appState.Env.VaultRoleID,
			SecretId: data2.(string),
		})
	if err != nil {
		logrus.Error(err.Error())
		return err
	}

	if err = client.SetToken(authAppRole.Auth.ClientToken); err != nil {
		logrus.Fatal(err)
	}

	// Доступ к секретам
	for _, path := range pathConfig {
		secret, err := client.Secrets.KvV2Read(context.Background(), path.Secrets.Path, vault.WithMountPath(path.Secrets.Engine))
		if err != nil {
			logrus.Infof(path.Secrets.Path)
			logrus.Infof(path.Secrets.Engine)
			logrus.Error(err.Error())
			return err
		}
		// Извлекаем конкретное поле из данных секрета, используя `path.Secrets.Field`
		nameValue, found := secret.Data.Data[path.Secrets.Field]
		if !found {
			return fmt.Errorf("поле %s не найдено в secret.Data", path.Secrets.Field)
		}

		err = setConfigServiceField(nameValue, path.Secrets.Name, appState.ConfigService)
		if err != nil {
			logrus.Error(err.Error())
			return err
		}
	}
	return nil
}

func setConfigServiceField(value interface{}, fieldName string, configService *models.ConfigService) error {
	configValue := reflect.ValueOf(configService).Elem()
	return setFieldRecursive(configValue, fieldName, value)
}

// Рекурсивная функция для установки значения, поддерживающая вложенные структуры
func setFieldRecursive(v reflect.Value, fieldName string, value interface{}) error {
	v = reflect.Indirect(v)
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("ожидалась структура, но получен %s", v.Kind())
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		jsonTag := field.Tag.Get("json")

		if jsonTag == fieldName && v.Field(i).CanSet() {
			fieldValue := v.Field(i)
			// Проверяем, что типы совпадают и устанавливаем значение
			if reflect.TypeOf(value).AssignableTo(fieldValue.Type()) {
				fieldValue.Set(reflect.ValueOf(value))
				return nil // Поле установлено, выходим из функции
			} else {
				return fmt.Errorf("несовместимый тип для поля %s", fieldName)
			}
		}
		if field.Type.Kind() == reflect.Struct {
			err := setFieldRecursive(v.Field(i), fieldName, value)
			if err == nil {
				return nil
			}
		}
	}

	return fmt.Errorf("поле с тегом %s не найдено в структуре ConfigService", fieldName)
}

func readVaultFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func loadEnvironment(Env *models.Environment) error {
	if err := godotenv.Load(envFilePath); err != nil {
		logrus.Warning("load file not found, environment variables load from environment")
	}
	if err := env.Parse(Env); err != nil {
		return err
	}
	return nil
}

func loadConfig(Config *models.ConfigService) error {
	file, err := os.ReadFile(configFilePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, Config)
	if err != nil {
		return err
	}

	return nil
}

func (f *customFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	prefixPath := pwd + "/"

	shortFilePath := strings.TrimPrefix(filepath.ToSlash(entry.Caller.File), filepath.ToSlash(prefixPath))

	var fields string
	for key, value := range entry.Data {
		fields += fmt.Sprintf("\"%s\":\"%v\",", key, value)
	}

	if len(fields) > 0 {
		fields = fields[:len(fields)-1]
	}

	if len(fields) > 0 {
		fields = ", " + fields
	}

	log := fmt.Sprintf(
		"{\"level\":\"%s\",\"msg\":\"%s\",\"point\": \" %s:%d \",\"short_point\":\"%s:%d\", \"time\":\"%s\"%s}\n",
		entry.Level.String(),
		entry.Message,
		entry.Caller.File,
		entry.Caller.Line,
		shortFilePath,
		entry.Caller.Line,
		entry.Time.Format(time.RFC3339),
		fields,
	)
	return []byte(log), nil
}

func RunLogger() error {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&customFormatter{})

	currentTime := time.Now()
	yearMonthDir := fmt.Sprintf("logs/%d-%02d", currentTime.Year(), currentTime.Month())

	err := os.MkdirAll(yearMonthDir, os.ModePerm)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}

	logFile := fmt.Sprintf("%s/%d-%02d-%02d.log", yearMonthDir, currentTime.Year(), currentTime.Month(), currentTime.Day())

	logFileHandle, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}

	logrus.SetOutput(io.MultiWriter(os.Stdout, logFileHandle))

	return nil
}
