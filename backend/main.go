package main

import (
	web "veles/webcontrol"
)

func main() {
	web.InitRouting("../keys/key.json", "physicsbridge-f4819.appspot.com")
}
