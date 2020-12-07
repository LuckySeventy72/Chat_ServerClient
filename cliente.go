 
package main

import (
	"fmt"
	"net"
	"time"
	"os"
	"bufio"
	"encoding/gob"
	"io/ioutil"
	"strings"
	"./clear"
	"path/filepath"

)

type Chat struct {
	
	Text []string
	ArchivoN string //nombre del archivo
	Archivo bool // si es true es text de archivo, si falso es mensaje para el chat
	Persona string
	Estado bool//Connect bool // si es true, el sujeto esta conecta al chat, si falso el sujeto se deconecto
}


func cliente() {
	var chat Chat
	c, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println("err1")
		fmt.Println(err)
		return
	}
	tempslChat:= []string{}
	//temp:= []string{}
	defer c.Close()
	//fmt.Print("cls")
	nickn:= ""
	fmt.Println(nickn)
	fmt.Println("Nickname: ")
	fmt.Scanln(&nickn)
	t:="Se conecto " + nickn
	chat.Text = append(chat.Text, t)
	chat.Archivo= false
	chat.Persona= nickn//nick name, nombre de usuario
	chat.Estado= true
	err2:= gob.NewEncoder(c).Encode(chat)//make it send slice
	if err2 != nil {
		fmt.Println(err2)
	//	return
	} else{
		fmt.Println("Conectando....")
		time.Sleep(time.Millisecond * 1000)
	}
	clear.CallClear()

	
	tempslChat = append(tempslChat,"null" )
	for{
		
		//cont:=0
		var chat Chat
		index:=0
		
		cha:=make (chan []string)
		cont:=0
		//lenChat:=0
		go func(){
			
			slChat:= []string{}
			
			for{
				time.Sleep(time.Millisecond * 500)
				err3:=gob.NewDecoder(c).Decode(&slChat)// recibir mensajes del servidor
				if err3 != nil {				
					continue
				
				}else{
									
					//if cont < len(slChat) {  
					for  cont <len(slChat) {
						if tempslChat[cont] != slChat[cont]{ 
							tempslChat[cont] = (slChat[cont] )
							tempslChat = append(tempslChat, "null" )
							cha<- tempslChat
						
						}
						cont++
					}						
				} 
			}
		}()
		
		go func(){
			
			for{
				tt:=<-cha
				index=0
				clear.CallClear()
				fmt.Println("---- ", nickn," ----")	
				fmt.Println("1. Enviar Mensaje")
				fmt.Println("2. Enviar Archivo")
				fmt.Println("0. Salir")
				l:= (len(tt) - 1)
				//time.Sleep(time.Millisecond * 500)
				for index < l{
				//	time.Sleep(time.Millisecond * 500)
					fmt.Println(tt[index])
					index=index +1

				}				
			}
		}()		
		var op string
		fmt.Scanln(&op)
		switch (op){
		case "1":{
			scanner := bufio.NewScanner(os.Stdin)//leer entrada, inluyendo espacios y saltos de linea
			fmt.Print("Mensaje:")
			if scanner.Scan() {
				texto := scanner.Text()
				chat.Text = append(chat.Text, nickn + ": " + texto)
				chat.Archivo= false
				chat.Persona=""
				chat.Estado= true
				//cont++
				gob.NewEncoder(c).Encode(chat)
				continue
				
			}
		}
		case "2": {// mandar archivo
			var dir string
			fmt.Println("Ingrese direccion del archivo:")
			fmt.Scanln(&dir)
			fileBytes, err := ioutil.ReadFile(dir)
			if err !=nil{
				fmt.Println(err)
			}else{
				//C:/Users/idont/Documents/proyecto_final/test04.txt
				path := dir
				fileN:= filepath.Base(path)//nombre del puro archivo
				chat.Text= strings.Split(string(fileBytes), "\n\r")
				//fmt.Println(chat.Text)
				chat.ArchivoN = fileN
				chat.Archivo = true
				chat.Persona=nickn
				chat.Estado= true
				err5:= gob.NewEncoder(c).Encode(chat)
				if err5 != nil {
				//	fmt.Println(err5)
				}
			}
			
		} 
		case "0": {		
				fmt.Println("bye, " + nickn)	
				chat.Text = append(chat.Text,"Se deconecto " + nickn)
				chat.Archivo= false
				chat.Persona= nickn
				chat.Estado= false
				gob.NewEncoder(c).Encode(chat)
				time.Sleep(time.Millisecond * 1000)
				//c.Close()
				return
			
		}
		}
	}//end of main loop
}

func main()  {
	cliente()	
	fmt.Println("fin")
}