package article_repo

import (
	"context"
	"edjr-trk/configs/env"
	"edjr-trk/configs/mongo"
	"edjr-trk/internal/repository"
	"edjr-trk/pkg/log"
	"fmt"
	"github.com/stretchr/testify/assert"
	m "go.mongodb.org/mongo-driver/mongo"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Инициализируем логгер и MongoDB

	// Запуск тестов
	os.Exit(m.Run())
}

func setup() (context.Context, repository.ArticleRepositoryInterface) {
	// Загружаем переменные окружения
	env.LoadEnv()

	// Инициализируем логгер и клиент MongoDB
	log.InitLogger()
	mongo.InitMongoSingleton()
	logger := log.GetLogger()
	clientDB := mongo.GetClient()

	// Создаем контекст
	ctx := context.Background()

	// Создаем репозиторий
	repo := repository.NewArticleRepository(clientDB, logger)

	return ctx, repo
}

func TestGetArticleById(t *testing.T) {
	ctx, repo := setup()
	// Подготовка тестовых данных
	validID := "674b0981fd898a8a128c5ffb"
	invalidID := "674b0981fd898a8a128c5fff"

	// Запуск теста
	t.Run("Success", func(t *testing.T) {
		article, err := repo.GetArticleById(ctx, validID)
		// Вывод структуры с ключами
		fmt.Printf("Article: %+v\n", article)

		// Проверяем, что ошибки нет
		assert.NoError(t, err)

		// Проверяем, что статья не nil
		assert.NotNil(t, article)

		// Проверяем, что структура содержит нужные ключи
		assert.Equal(t, "text 1732970881274", article.Text)
		assert.Equal(t, "Hi", article.Title)
		assert.NotEmpty(t, article.Img)
		assert.NotEmpty(t, article.Date)
		assert.NotEmpty(t, article.ID)
	})

	// Запуск теста
	t.Run("Not found", func(t *testing.T) {
		article, err := repo.GetArticleById(ctx, invalidID)
		// Вывод структуры с ключами
		fmt.Printf("Article: %+v\n", article)
		fmt.Printf("err: %+v\n", err)
		// Проверяем, что вернулась ошибка mongo.ErrNoDocuments
		assert.ErrorIs(t, err, m.ErrNoDocuments)
	})
}
