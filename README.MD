# WasmPulse


***WasmPulse*** WasmPulse is an open-source WebAssembly (Wasm) monitoring tool, created to fill in the gaps in monitoring Server-side Wasm services. It monitors for a PID running a Wasm runtime, then starts tracking it’s resource (CPU, Memory) usage from the OS’s perspective. It exposes an API with these metrics in a Prometheus-compatible format.
Written in Golang, the WasmPulse tool also features testing scripts for running sample Wasm services, and a Docker compose to facilitate bringing up WasmPulse, Prometheus and Grafana, all pre-configured with 1 command.


## TL;DR
```console
cd release
docker compose up --build
```

## Documentation
The documentation can be found [here](./docs/)!

## Contributing
MInterested in contributing? Check out [contributions](./CONTRIBUTING.md).

## License
***WasmPulse*** is under the [MIT](./LICENSE) license.

## Citation
Text format:
```
M. Ziółkowski, I. Fajjari and H. Tarasiuk, "WebAssembly as a Lightweight Path to Sustainable and High-Performance Cloud-Native Computing," 2026 IEEE 23rd Consumer Communications & Networking Conference (CCNC), Las Vegas, NV, USA, 2026, pp. 1-6, doi: 10.1109/CCNC65079.2026.11366292.
```
BibTex:
```
@INPROCEEDINGS{11366292,
  author={Ziółkowski, Marcin and Fajjari, Ilhem and Tarasiuk, Halina},
  booktitle={2026 IEEE 23rd Consumer Communications & Networking Conference (CCNC)}, 
  title={WebAssembly as a Lightweight Path to Sustainable and High-Performance Cloud-Native Computing}, 
  year={2026},
  volume={},
  number={},
  pages={1-6},
  keywords={Cloud computing;Memory management;Containers;Webassembly;Telecommunications;Sustainable development;Virtualization;Environmental impact;Next generation networking;Load modeling;Virtualization;WebAssembly;Sustainability;Containerization;Cloud Native},
  doi={10.1109/CCNC65079.2026.11366292}}

```
