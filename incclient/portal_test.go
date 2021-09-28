package incclient

import (
	"fmt"
	"testing"
	"time"
)

var ICPortal *IncClient

func initICPortal() error {
	var err error
	ICPortal, err = NewTestNetClient()
	if err != nil {
		return fmt.Errorf("cannot init new incognito client")
	}

	return nil
}

func TestIncClient_CreateAndSendPortalShieldTransaction(t *testing.T) {
	// test cases
	type TestCasePortalShield struct {
		paymentAddr    string
		shieldingProof string
	}

	tcs := []TestCasePortalShield{
		{
			paymentAddr:    "12smNK6U7rRbxLConJrmjHGgYFhNVTmKNoqn8B8rJuk5J2ZY363yCsSdAmrbnhrMtNHuXzszRB1xX8VGe6FuxjVqJWwhMmxDKuoaGZfuUaLAC2qnozu2czneFyvTUVAh4kaqLft1yEe5jRydnh39",
			shieldingProof: "eyJNZXJrbGVQcm9vZnMiOlt7IlByb29mSGFzaCI6Wzc4LDI1NSwyMzEsMTcyLDUwLDI1MiwyMSwyMjMsNjcsMTU2LDYxLDE5NiwyMTQsODMsMTg5LDIwLDE2MCwxMzcsMjA4LDgwLDEzNiwyOSw3OSw1LDMsNzIsMTU1LDE3Nyw5NywxMTksMjE5LDE1NV0sIklzTGVmdCI6dHJ1ZX0seyJQcm9vZkhhc2giOlsyNDMsMzUsMTE0LDg4LDExNiwxMTksNjcsODIsNTksMjU1LDU3LDE4NSwxNjksMjQxLDEzNiwyOSwzMSwyMTAsMTIzLDIzNCwxOTcsNTQsNDksMzUsNTksMTgxLDEwMSw4NSw1MCwxNjcsMjgsMjA1XSwiSXNMZWZ0Ijp0cnVlfSx7IlByb29mSGFzaCI6WzEzMiw5Myw2NSw4LDEsMjI1LDU0LDgsMTI2LDE3LDEwMywxMTksMTQ5LDE0NSwyNDMsMjI1LDI5LDE0NiwzMywxNiwxNjIsMjQ1LDMwLDEyLDE5MiwxNzMsMTA0LDUzLDQ4LDE5NywxOTUsNjVdLCJJc0xlZnQiOnRydWV9LHsiUHJvb2ZIYXNoIjpbOTksMTQ1LDE0MywxMiwzOSw5MSwyMDQsMTcwLDk2LDE2Nyw2MSwyMjQsMjQwLDI0MSwxMDMsMTgzLDY0LDczLDUsNDcsMjI1LDE0Miw3MSwxNDMsMjI0LDMyLDQ5LDI0NSw4NSwxNTUsMjIyLDQzXSwiSXNMZWZ0Ijp0cnVlfSx7IlByb29mSGFzaCI6WzQzLDIwMiwxNTEsNTcsMTMxLDIyMCwzOCw5MSw0MCwxNzYsMTkwLDIxNiwxNywxNjIsMTEzLDQ2LDE3NiwxNzAsNDQsMjI2LDE1OSw2MiwxODQsNTcsMTQ4LDU5LDE3Miw2OCwxMywxOTMsMiwxMjFdLCJJc0xlZnQiOmZhbHNlfSx7IlByb29mSGFzaCI6WzE1NCwxNDUsMjEzLDE4OSwxOCwzMSw5MSw1MCw3LDgwLDE3NSwxNzAsNjAsOTYsMTIwLDEyMywxNTYsMTI4LDE4Miw3LDEzLDE5MSw4NCwxOTQsMTM3LDE3MiwxMjYsMTgwLDAsMjM3LDE4OCw5MV0sIklzTGVmdCI6dHJ1ZX1dLCJCVENUeCI6eyJWZXJzaW9uIjoyLCJUeEluIjpbeyJQcmV2aW91c091dFBvaW50Ijp7Ikhhc2giOlsxMjMsNzEsMjExLDIxMSwyMTgsMjI2LDEzNiwxNjYsMTcxLDUyLDcxLDMxLDE0MiwxOTQsMjUwLDU5LDI5LDc4LDU1LDExMiwxMjEsNTIsNTksNDcsMjQ2LDksMTQ2LDcxLDE0MCw4LDE3MywxNzhdLCJJbmRleCI6MX0sIlNpZ25hdHVyZVNjcmlwdCI6IiIsIldpdG5lc3MiOm51bGwsIlNlcXVlbmNlIjo0Mjk0OTY3MjkzfV0sIlR4T3V0IjpbeyJWYWx1ZSI6MTAwMDAwLCJQa1NjcmlwdCI6IkFDQ1dnclFYekZaamFUQUZ0YnlIY2lMRkhuWmp1WFZUZW5RYTBmeXd2c094anc9PSJ9LHsiVmFsdWUiOjEzMTAwODgsIlBrU2NyaXB0IjoiQUJTdkpSM3dFQjJ4TklEcTkyV1lHTFBPSk8xM0VBPT0ifV0sIkxvY2tUaW1lIjoyMDk2NjU3fSwiQmxvY2tIYXNoIjpbMSw5NSw3MCw3OSwxMjMsNDEsNzIsMTQ5LDE1Nyw1OSwyMjEsMTk5LDksMTE4LDEwMSw4MywyNDYsMTcsMjUyLDEzNCwyMjUsODksMTA0LDg5LDEsMCwwLDAsMCwwLDAsMF19",
		},
		{
			paymentAddr:    "12sw5DqMY42zcVbnwxHpGWmVykQ1W89UJexJug62qUK7GWi9x88vtet1ZM4M5tmkWec4tXfeURV9WAFgtAwHs1AitaSW1iVUiX8m4HG3m3NnuV74SYcbJBNFFZPb6PrVvDk9wZmGzSUKSd2xRd4L",
			shieldingProof: "eyJNZXJrbGVQcm9vZnMiOlt7IlByb29mSGFzaCI6WzE2MCwxODIsMTEyLDEzMywyMzcsODYsMjA3LDIxNSw1MiwyMDksMjEsMzQsMjE4LDU0LDEwOCwxNjcsMTc1LDMwLDE4MywxMTEsMTY1LDI0NCwxMjMsNTksODcsMjIzLDExMSwxMzUsMjI5LDEzMyw2MCwxNzJdLCJJc0xlZnQiOnRydWV9LHsiUHJvb2ZIYXNoIjpbNDQsMTk1LDIzMCwxNCwyMDgsMTcxLDIwLDQyLDcsNDksNjMsMTM3LDEyMywzMywxNywyMCwyMTgsNDcsMTY2LDExMiwxNzAsNjUsMiwxMzYsMTA4LDIwMywyMjksMTc4LDIxOCw0OSwzNiw4MV0sIklzTGVmdCI6ZmFsc2V9LHsiUHJvb2ZIYXNoIjpbMTEzLDIxNywyMDQsMTQzLDEwMSwxODIsMTA1LDI1NSw5OCw2NCwxNzYsOTUsMjUxLDY4LDI0NSwyMDEsMTE4LDE3NywxNTcsMTE5LDE5MSwyMTEsMjQzLDEzNywyMDMsMTg3LDEzMiwxNjIsMTg1LDY1LDEzNSwyMzZdLCJJc0xlZnQiOmZhbHNlfSx7IlByb29mSGFzaCI6WzIyNSwxNiwxMjcsMTQyLDIyNiw1OCwxMTYsNywxNDAsMjA4LDIwNiwxNjAsMTc3LDM2LDE2MCw5NSwxNjMsMjQsMjEsNDAsMTksMTMsMTQzLDExNiwxODMsMTQzLDIzNywxNjUsODQsNDQsOTMsMTNdLCJJc0xlZnQiOmZhbHNlfSx7IlByb29mSGFzaCI6WzQzLDIwMiwxNTEsNTcsMTMxLDIyMCwzOCw5MSw0MCwxNzYsMTkwLDIxNiwxNywxNjIsMTEzLDQ2LDE3NiwxNzAsNDQsMjI2LDE1OSw2MiwxODQsNTcsMTQ4LDU5LDE3Miw2OCwxMywxOTMsMiwxMjFdLCJJc0xlZnQiOmZhbHNlfSx7IlByb29mSGFzaCI6WzE1NCwxNDUsMjEzLDE4OSwxOCwzMSw5MSw1MCw3LDgwLDE3NSwxNzAsNjAsOTYsMTIwLDEyMywxNTYsMTI4LDE4Miw3LDEzLDE5MSw4NCwxOTQsMTM3LDE3MiwxMjYsMTgwLDAsMjM3LDE4OCw5MV0sIklzTGVmdCI6dHJ1ZX1dLCJCVENUeCI6eyJWZXJzaW9uIjoyLCJUeEluIjpbeyJQcmV2aW91c091dFBvaW50Ijp7Ikhhc2giOls1NywxMDEsOTQsODUsMzMsNzUsNTksMTgzLDIxLDc4LDcxLDEzNCwxNzcsMTU2LDIwLDE3MywxNTQsNDgsMTA4LDUsODgsMjE4LDE2NSwyMTEsMjUzLDgsMTg4LDk1LDE2MSw2MCw4MiwyMjFdLCJJbmRleCI6MH0sIlNpZ25hdHVyZVNjcmlwdCI6IiIsIldpdG5lc3MiOm51bGwsIlNlcXVlbmNlIjo0Mjk0OTY3MjkzfSx7IlByZXZpb3VzT3V0UG9pbnQiOnsiSGFzaCI6WzY3LDQ3LDI1LDc5LDE4OSwxMTksOCwxMjgsMjAwLDEwOCw2LDQ1LDMxLDQwLDIwOCwxNzcsMTYyLDE1OCwxMzQsMTQsMTM3LDM3LDM5LDE3MywxMDIsMjEyLDQ3LDYwLDIyMSw5OSwxMzMsMTRdLCJJbmRleCI6MH0sIlNpZ25hdHVyZVNjcmlwdCI6IiIsIldpdG5lc3MiOm51bGwsIlNlcXVlbmNlIjo0Mjk0OTY3MjkzfV0sIlR4T3V0IjpbeyJWYWx1ZSI6MTU3NjAsIlBrU2NyaXB0IjoiQUJRa3MyVVAvUlRqcjhOakFHZ0RjQktiMEYwdE9nPT0ifSx7IlZhbHVlIjoxNTAwMDAsIlBrU2NyaXB0IjoiQUNENitOeUZoOWNsWUt5ZmlHMTl1QjhodGVjRXltQWhlQ25CRmpHNGVwYnhYQT09In1dLCJMb2NrVGltZSI6MjA5NjY1OH0sIkJsb2NrSGFzaCI6WzEsOTUsNzAsNzksMTIzLDQxLDcyLDE0OSwxNTcsNTksMjIxLDE5OSw5LDExOCwxMDEsODMsMjQ2LDE3LDI1MiwxMzQsMjI1LDg5LDEwNCw4OSwxLDAsMCwwLDAsMCwwLDBdfQ==",
		},
	}

	err := initICPortal()
	if err != nil {
		panic(err)
	}

	// Input your private key - to pay transaction fee
	privateKey := ""
	tokenID := "4584d5e9b2fc0337dfb17f4b5bb025e5b82c38cfa4f54e8a3d4fcdd03954ff82"

	for _, tc := range tcs {
		shieldID, err := ICPortal.CreateAndSendPortalShieldTransaction(privateKey, tokenID, tc.paymentAddr,
			tc.shieldingProof, nil, nil)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Shield request sent and waiting for complete\n")
		time.Sleep(50 * time.Second)

		fmt.Printf("Check shielding status\n")
		status, err := ICPortal.GetPortalShieldingRequestStatus(shieldID)
		if err != nil {
			panic(err)
		}
		if status.Status == 1 {
			fmt.Printf("Shield completed\n")
		} else {
			fmt.Printf("Shield reject with error %v\n", status.Error)
		}
	}
}

func TestIncClient_CreateAndSendPortalUnShieldTransaction(t *testing.T) {
	// test cases
	type TestCasePortalUnShield struct {
		unshieldAmount uint64
		remoteAddress  string
	}

	tcs := []TestCasePortalUnShield{
		{
			unshieldAmount: 0.0005 * 1e9,
			remoteAddress:  "tb1q0qjpqrgz54xsseymjrql6xs7p9qm0uj53v6rw9",
		},
		{
			unshieldAmount: 0.0005 * 1e9,
			remoteAddress:  "tb1q0qjpqrgz54xsseymjrql6xs7p9qm0uj53v6rw9",
		},
	}

	err := initICPortal()
	if err != nil {
		panic(err)
	}

	// Input your private key - to pay transaction fee
	privateKey := ""
	tokenID := "4584d5e9b2fc0337dfb17f4b5bb025e5b82c38cfa4f54e8a3d4fcdd03954ff82"

	for _, tc := range tcs {
		unshieldID, err := ICPortal.CreateAndSendPortalUnShieldTransaction(privateKey, tokenID, tc.remoteAddress,
			tc.unshieldAmount, nil, nil)
		if err != nil {
			panic(err)
		}
		fmt.Printf("unshieldID: %v\n", unshieldID)
		fmt.Printf("UnShield request sent and waiting for complete\n")
		time.Sleep(40 * time.Second)

		fmt.Printf("Check unshielding status\n")

		status, err := ICPortal.GetPortalUnShieldingRequestStatus(unshieldID)
		if err != nil {
			panic(err)
		}
		if status.Status == 0 || status.Status == 1 {
			fmt.Printf("Unshield request is processing\n")
		} else if status.Status == 3 {
			fmt.Printf("Unshield request is rejected\n")
		}
	}
}
