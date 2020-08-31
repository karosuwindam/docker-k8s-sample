package main

import (
	"fmt"
)

func main() {
	cpu_tmp := cpuTmpOutput(useIoutilReadFile(CPU_TMP_PASS))                 //CPU temp
	gpu_tmp := tmpOut(cmdOutput("vcgencmd", "measure_temp"))                 //GPU temp
	clock := clockOutput(cmdOutput("vcgencmd", "measure_clock", "arm"))      //ARM cores
	vc4_core := voltOutput(cmdOutput("vcgencmd", "measure_volts", "core"))   //VC4 core voltage
	sdram_c := voltOutput(cmdOutput("vcgencmd", "measure_volts", "sdram_c")) //SDRAM Core Voltage
	sdram_i := voltOutput(cmdOutput("vcgencmd", "measure_volts", "sdram_i")) //SDRAM I/O voltage
	sdram_p := voltOutput(cmdOutput("vcgencmd", "measure_volts", "sdram_p")) //SDRAM Phy Voltage
	mem_arm := memOut(cmdOutput("vcgencmd", "get_mem", "arm"))               //Mem arm
	mem_gpu := memOut(cmdOutput("vcgencmd", "get_mem", "gpu"))               //Mem gpu
	fmt.Println(cpu_tmp)                                                     //CPU temp
	fmt.Println(gpu_tmp)
	fmt.Println(clock)
	fmt.Println(vc4_core)
	fmt.Println(sdram_c)
	fmt.Println(sdram_i)
	fmt.Println(sdram_p)
	fmt.Println(mem_arm)
	fmt.Println(mem_gpu)
}
