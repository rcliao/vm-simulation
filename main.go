package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
)

// VM is a Simulated VM model to contain common VM attributes
type VM struct {
	MaxMemory, MaxStorage         int32
	CPU, Memory, Storage, Network []int32
}

var randomVMs []VM

func init() {
	rand.Seed(594)
	for i := 0; i < 1000; i++ {
		maxMemory := rand.Int31n(80000) + 1
		maxStorage := rand.Int31n(80000) + 1
		var cpuHistory, memoryHistory, storageHistory, networkHistory []int32
		CPU := rand.Int31n(100)
		Memory := rand.Int31n(maxMemory)
		Storage := rand.Int31n(maxStorage)

		for j := 0; j < 1000; j++ {
			switch i {
			case 437:
				networkHistory = append(networkHistory, 0)
				break
			case 454:
				CPU = continueRandomChange(CPU, 100, 20)
				Memory = continueRandomChange(Memory, maxMemory, maxMemory/5+10)
				Storage = continueRandomChange(Storage, maxStorage, maxStorage/5+10)
				networkHistory = append(networkHistory, rand.Int31n(100))
				break
			default:
				CPU = continueRandomChange(CPU, 100, 10)
				Memory = continueRandomChange(Memory, maxMemory, maxMemory/100+1)
				Storage = continueRandomChange(Storage, maxStorage, maxStorage/100+1)
				networkHistory = append(networkHistory, rand.Int31n(100))
				break
			}
			cpuHistory = append(cpuHistory, CPU)
			memoryHistory = append(memoryHistory, Memory)
			storageHistory = append(storageHistory, Storage)
		}
		randomVMs = append(randomVMs, VM{
			maxMemory,
			maxStorage,
			cpuHistory,
			memoryHistory,
			storageHistory,
			networkHistory,
		})
	}
}

func continueRandomChange(initialValue, max, variance int32) int32 {
	cpuChange := rand.Int31n(variance)
	result := initialValue
	if rand.Float32() > 0.5 {
		if result+cpuChange > max {
			result -= cpuChange
		} else {
			result += cpuChange
		}
	} else {
		if result-cpuChange < 0 {
			result += cpuChange
		} else {
			result -= cpuChange
		}
	}
	return result
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/api/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!\n")
	})

	e.GET("/api/vm/:id", func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		return c.JSON(http.StatusOK, randomVMs[id])
	})

	fmt.Println("Server running at port 9000")

	e.Run(standard.New(":9000"))
}
