package main

import (
	"fmt"
	"net"
	"encoding/gob"
	"time"
	"os"
	"strings"
	"./clear"
	
)

type Chat struct {
	Text []string
	Persona string
	ArchivoN string //nombre del archivo
	Archivo bool// si true es text de archivo, si falso es text del chat
	Estado bool// si es true, el sujeto esta conecta al chat, si falso el sujeto se deconecto
}
func servidor()  {
	
	
	s, err := net.Listen("tcp", ":9999")
	if err != nil {
		
		fmt.Println(err)
		
		return
	}
	defer s.Close()
	var chat Chat
	slChat:= []string{}
	slFileText:= []string{}//read	
	slFileName:= []string{}//read	
	users:= []string{}//personas activo en el chat
	//contM:=0
	go func(){
			for{
				time.Sleep(time.Millisecond * 500)
				clear.CallClear()
				fmt.Println("1. Mostrar Chat")
				fmt.Println("2. Mostrar Archivos")
				fmt.Println("3. Respaldar Archivo/Mensajes")
				fmt.Println("0. Salir")
				var op int
				fmt.Scanln(&op)

				if op==1{
					contM:= 0//contador
					var inp string
					clear.CallClear()
					go func(){
						//for{
							//for _,a := range slChat{
							//	fmt.Println(a)
							//} 
						//}
						for{
							for contM < len(slChat){
								fmt.Println(slChat[contM])
								contM++
							} 
						}
					}()
					fmt.Scanln(&inp)

					//time.Sleep(time.Millisecond * 4000)
				}
				if op==2{
					var inp string
					clear.CallClear()
					go func(){
						fmt.Println("Nombres de Archivos y su Contenido:")
						for _,a:= range slFileName{
							fmt.Println(a)
						} 
						fmt.Println("\n-----------\n")
						for i,a := range slFileText{
							fmt.Println(slFileName[i], ":", a)
						} 
					}()
					fmt.Scanln(&inp)

					//time.Sleep(time.Millisecond * 4000)
				}
				if op==3{
					for i,fn:= range slFileName{//respalde de archivos

						nf, err:=os.Create(("Respaldo" + fn))
						if err !=nil{
						
							fmt.Println(err)
						}
						tx:= slFileText[i]
						nf.WriteString(tx)
						nf.Close()
					}

					nf2, err:=os.Create(("RespaldoChat.txt"))
					if err !=nil{
						
						fmt.Println(err)
					}
				
					for _,sy:= range slChat{
						nf2.WriteString(sy+ string('\n'))
					}


				}
				if op==0{
					os.Exit(0)
				}	
		}
	}()
	for {
		c, err := s.Accept()
		if err != nil {

			//fmt.Println(err)
			continue
		}

		go handleClient(c, &slChat, &slFileName, &slFileText, &chat, &users)	
	}
}

func handleClient( c net.Conn, slChat* []string, slFileName* []string, slFileText* []string, chat* Chat, users* []string){

	go func(){

		for {
			*chat =Chat{}
			time.Sleep(time.Millisecond * 500)
			err3:= gob.NewDecoder(c).Decode(&chat)
			if err3 != nil {
				continue
			}else{ 
					if chat.Estado ==false{ // se elimina usuario de lista cuando se desconecta
						index:=0
						for y, name:= range *users{
							if name == chat.Persona{
								index=y
								break
							}
						}
						*users = removeIndex(*users, index)
						chatString:= strings.Join(chat.Text," ") //convierte un slice en un string
						*slChat = append(*slChat, chatString)
						cha:=make (chan []string)
						go actualizar(cha, *slChat)
						//c.Close()
						continue
					}
					
					if chat.Archivo != true{// para mensajes del chat, se agrega mensajes en un slice
						*users = append(*users, chat.Persona)
						chatString:= strings.Join(chat.Text," ")
						*slChat = append(*slChat, chatString)
						cha:=make (chan []string)
						go actualizar(cha, *slChat)

					}else{//para manejo del archivo, se agrega el texto en un slice
						
						*slFileName = append(*slFileName, chat.ArchivoN)
						ArchivoString:= strings.Join(chat.Text,"\n\r ")
						*slFileText = append(*slFileText, ArchivoString) //almacena el texto del archivo en un slice
						*slChat = append(*slChat, chat.Persona + " envio un archivo, " + chat.ArchivoN)
						
						chatString:= strings.Join(chat.Text," ")
						cha:=make (chan []string)
						*slChat = append(*slChat, "     CONTENIDO\n"+ "------------------ \n"+   chatString + "\n------------------")
						go actualizar(cha, *slChat)
						continue

					}

				}

			}
		}()

	go func(){///--enviar chat a clientes
		for{
			time.Sleep(time.Millisecond * 500)
			cha:=make (chan []string)
			go actualizar(cha, *slChat)
			chanChat:=<- cha

			err4:= gob.NewEncoder(c).Encode(chanChat)
			if err4 != nil {
				continue		
			}
		}
	}()
}

func actualizar(ch chan []string, chChat []string){//actualizar Slice de los mensajes del chat
	for{
	ch<-chChat
	}  
}

func removeIndex(s []string, index int) []string {//eliminar a persona de slice se desconecte
    return append(s[:index], s[index+1:]...)
}

func main() {

	servidor()
	var input string
	fmt.Scanln(&input)
	
}