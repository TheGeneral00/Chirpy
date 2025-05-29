package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

// testing auth jwt functions with random uuids

func TestAuthJWT(t *testing.T){
        uuids := []uuid.UUID{}
        tokenSecret := "myTestSecret"
        expiresIn, err := time.ParseDuration("300s")
        if err != nil {
                t.Errorf("Not able to parse duration\n")
        }
        i := 0
        for i < 10 {
                newUUID, err := uuid.NewUUID()
                if err != nil {
                        t.Errorf("Failed to generate uuid no. %v\n", i)
                        continue
                }
                uuids = append(uuids, newUUID)
                i++
        }
        
        for _, uuid := range(uuids){
                jwt, err := MakeJWT(uuid, tokenSecret, expiresIn)
                if err != nil {
                        t.Errorf("Failed to make jwt: %v\n", err)
                }
                rediscoveredUuid, err := ValidateJWT(jwt, tokenSecret)
                if err != nil {
                        t.Errorf("Failed to parse jwt: %v\n", err)
                }
                if uuid != rediscoveredUuid {
                        t.Errorf("UUID not equal: %v %v\n", uuid, rediscoveredUuid)
                }
        }
}
