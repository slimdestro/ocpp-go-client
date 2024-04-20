# 
ocpp-go-client

## ocpp-go-client
#### Golang Client package for Open Charge Point Protocol (OCPP)

[![N|Solid](https://miro.medium.com/v2/resize:fit:1200/1*J-HlYNxNOE3ilmPgxv0_Ug.png)](https://www.modcode.dev)
 
## Installation

Configuration

```sh
import "github.com/slimdestro/ocpp-go-client"

endpoint := "ocpp-server-here"
client := ocpp.NewClient(endpoint)

chargeBoxID := "12345"
err := client.BootNotification(chargeBoxID)
if err != nil {
    fmt.Println("BootNotification failed:", err)
}

```
## Author

[slimdestro(Mukul Kumar)](https://linktr.ee/slimdestro)
