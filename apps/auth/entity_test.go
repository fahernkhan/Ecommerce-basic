package auth

import (
	"Ecommerce-basic/infra/response"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestValidateAuthEntity(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		authEntity := AuthEntity{
			Email:    "fathur@gmail.com",
			Password: "mysecretpassword",
		}

		err := authEntity.Validate()
		require.Nil(t, err)
	})

	t.Run("email is required", func(t *testing.T) {
		authEntity := AuthEntity{
			Email:    "",
			Password: "mysecretpassword",
		}

		err := authEntity.Validate()
		require.NotNil(t, err)
		require.Equal(t, response.ErrEmailRequired, err)
	})

	t.Run("email is invalid", func(t *testing.T) {
		authEntity := AuthEntity{
			Email:    "fathur.id",
			Password: "mysecretpassword",
		}

		err := authEntity.Validate()
		require.NotNil(t, err)
		require.Equal(t, response.ErrEmailInvalid, err)
	})
	t.Run("password is required", func(t *testing.T) {
		authEntity := AuthEntity{
			Email:    "reyhan@gmail.com",
			Password: "",
		}

		err := authEntity.Validate()
		require.NotNil(t, err)
		require.Equal(t, response.ErrPasswordRequired, err)
	})

	t.Run("password must have minimum 6 character", func(t *testing.T) {
		authEntity := AuthEntity{
			Email:    "reyhan@gmail.com",
			Password: "pass",
		}

		err := authEntity.Validate()
		require.NotNil(t, err)
		require.Equal(t, response.ErrPasswordInvalidLength, err)
	})
}

func TestEncryptPassword(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		authEntity := AuthEntity{
			Email:    "reyhan@gmail.com",
			Password: "mysecretpassword",
		}

		err := authEntity.EncryptPassword(bcrypt.DefaultCost)
		require.Nil(t, err)

		log.Printf("%+v\n", authEntity)
	})
}
