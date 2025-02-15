package main

import (
        "flag"
        "fmt"
        "net"
        "os"
)

func main() {
        // Определяем аргументы командной строки
        sourceIP := flag.String("source-ip", "0.0.0.0", "IP address to listen for incoming traffic")
        sourcePort := flag.Int("source-port", 19133, "Port to listen for incoming traffic")
        destIP := flag.String("dest-ip", "127.0.0.1", "IP address to forward traffic")
        destPort := flag.Int("dest-port", 19132, "Port to forward traffic")
        mode := flag.String("mode", "tcp-to-udp", "Mode of operation: udp-to-tcp or tcp-to-udp")

        // Парсим аргументы
        flag.Parse()

        // Настройки адресов
        sourceAddr := fmt.Sprintf("%s:%d", *sourceIP, *sourcePort)
        destAddr := fmt.Sprintf("%s:%d", *destIP, *destPort)

        switch *mode {
        case "udp-to-tcp":
                udpToTCP(sourceAddr, destAddr)
        case "tcp-to-udp":
                tcpToUDP(sourceAddr, destAddr)
        default:
                fmt.Println("Invalid mode. Use 'udp-to-tcp' or 'tcp-to-udp'.")
                os.Exit(1)
        }
}

func udpToTCP(sourceAddr, destAddr string) {
        // Создаем UDP-сокет
        udpAddr, err := net.ResolveUDPAddr("udp", sourceAddr)
        if err != nil {
                fmt.Println("Error resolving UDP address:", err)
                return
        }

        udpConn, err := net.ListenUDP("udp", udpAddr)
        if err != nil {
                fmt.Println("Error creating UDP socket:", err)
                return
        }
        defer udpConn.Close()

        fmt.Println("Listening for UDP traffic on", sourceAddr)

        buffer := make([]byte, 1024)

        for {
                n, addr, err := udpConn.ReadFromUDP(buffer)
                if err != nil {
                        fmt.Println("Error reading from UDP:", err)
                        continue
                }

                fmt.Printf("Received UDP data from %s: %s\n", addr.String(), buffer[:n])

                // Создаем TCP-сокет и отправляем данные
                tcpConn, err := net.Dial("tcp", destAddr)
                if err != nil {
                        fmt.Println("Error creating TCP connection:", err)
                        continue
                }

                _, err = tcpConn.Write(buffer[:n])
                if err != nil {
                        fmt.Println("Error writing to TCP:", err)
                }

                tcpConn.Close()
        }
}

func tcpToUDP(sourceAddr, destAddr string) {
        // Создаем TCP-сокет
        tcpListener, err := net.Listen("tcp", sourceAddr)
        if err != nil {
                fmt.Println("Error creating TCP socket:", err)
                return
        }
        defer tcpListener.Close()

        fmt.Println("Listening for TCP traffic on", sourceAddr)

        for {
                tcpConn, err := tcpListener.Accept()
                if err != nil {
                        fmt.Println("Error accepting TCP connection:", err)
                        continue
                }

                go handleTCPToUDP(tcpConn, destAddr)
        }
}

func handleTCPToUDP(tcpConn net.Conn, destAddr string) {
        defer tcpConn.Close()

        buffer := make([]byte, 1024)
        n, err := tcpConn.Read(buffer)
        if err != nil {
                fmt.Println("Error reading from TCP:", err)
                return
        }

        fmt.Printf("Received TCP data: %s\n", buffer[:n])

        // Создаем UDP-сокет и отправляем данные
        udpAddr, err := net.ResolveUDPAddr("udp", destAddr)
        if err != nil {
                fmt.Println("Error resolving UDP address:", err)
                return
        }

        udpConn, err := net.DialUDP("udp", nil, udpAddr)
        if err != nil {
                fmt.Println("Error creating UDP connection:", err)
                return
        }
        defer udpConn.Close()

        _, err = udpConn.Write(buffer[:n])
        if err != nil {
                fmt.Println("Error writing to UDP:", err)
        }
}
