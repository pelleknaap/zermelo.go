# zermelo.go
A little zermelo API wrapper, made in Go (golang).
I'll add a lot of other features in the near future but for now it only supports getting an access token with a "koppel code" and getting appointments.
Things I'm looking into adding are:
Function to get current announcements,


# Quick Start
Create a "ZermeloData" object, provide the access_token directly or use a "koppel code".
If you want to use a "koppel code" you need to call the GetApiKey function to get an access_token.
You can get all the appointments by calling the "GetAppointments" function.
Please look into the examples folder to see the code.

# Examples
You can find examples in the /examples folder. 

# Disclaimer
This project is not affiliated with Zermelo Software B.V in any way. It's only an api wrapper for Go.
