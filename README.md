https://github.com/rs/zerolog

nvidia-smi --id=0 --query-gpu=name --format=csv,noheader

echo; nvidia-smi --query-gpu=timestamp,name,pci.bus_id,driver_version,pstate,pcie.link.gen.max,pcie.link.gen.current,temperature.gpu,utilization.gpu,utilization.memory,memory.total,memory.free,memory.used --format=csv,noheader

[basics]
index
timestamp
name

[advanced]
pci.bus_id
driver_version
pstate - Power state
pcie.link.gen.max
pcie.link.gen.current

[charted]
temperature.gpu
utilization.gpu
utilization.memory
memory.total
memory.free
memory.used

2020/12/02 22:37:14.163, GeForce RTX 2080 Ti, 00000000:09:00.0, 450.80.02, P8, 3, 1, 39, 1 %, 3 %, 10997 MiB, 10200 MiB, 797 MiB


--- Explanation ---
timestamp	The timestamp of where the query was made in format "YYYY/MM/DD HH:MM:SS.msec".

name	The official product name of the GPU. 
This is an alphanumeric string. For all products.

pci.bus_id	PCI bus id as "domain:bus:device.function", in hex.
driver_version	The version of the installed NVIDIA display driver. 
This is an alphanumeric string.

pstate	The current performance state for the GPU. States range from P0 (maximum performance) to P12 (minimum performance).

pcie.link.gen.max	The maximum PCI-E link generation possible with this GPU and system configuration. 
For example, if the GPU supports a higher PCIe generation than the system supports then this reports the system PCIe generation.

pcie.link.gen.current	The current PCI-E link generation. These may be reduced when the GPU is not in use.
temperature.gpu	Core GPU temperature. in degrees C.

utilization.gpu
Percent of time over the past sample period during which one or more kernels was executing on the GPU.
The sample period may be between 1 second and 1/6 second depending on the product.

utilization.memory
Percent of time over the past sample period during which global (device) memory was being read or written.
The sample period may be between 1 second and 1/6 second depending on the product.

memory.total
Total installed GPU memory.

memory.free
Total free memory.

memory.used
Total memory allocated by active contexts.