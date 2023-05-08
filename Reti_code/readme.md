<!-- @format -->

# Reticulum: A Two-Layer Blockchain Sharding Protocol

### V 0.1.2

#### Note!

**This version of the code is currently being developed and improved. It is specifically designed for conducting comparative experiments and should not be used in real commercial applications as it may have some imperfections. Its purpose is to gather insights and data for further enhancements and optimizations before considering it for production environments.**

**If you aim to replicate the results of the experiments, we suggest utilizing a high-performance server to accommodate a larger number of nodes. However, if your intention is to gain a basic understanding of Reticulum's functionality, deploying around ten nodes would be a suitable option. Keep in mind that the specific hardware and server specifications can impact the performance and scalability of the system. It's important to consider these factors when setting up your parameters.**

- Commands for configuring the environment：

  ```
  chmod +x start.sh
  sudo add-apt-repository ppa:longsleep/golang-backports
  sudo apt update
  sudo apt install golang-1.19
  echo 'export PATH="/usr/lib/go-1.19/bin:$PATH"' >> ~/.bashrc
  source ~/.bashrc
  go version
  go mod init example.com/m
  go mod init example.com/m/v2
  go get gopkg.in/yaml.v2
  ```

- Network setting:

  ```
  sudo tc qdisc add dev lo root handle 1: tbf rate 10gbit burst 100000 latency 20ms sudo tc qdisc add dev lo parent 1:1 handle 10:

  sudo tc qdisc del dev lo root netem
  ```

- Parameters of reticulum:

  You have the flexibility to set the parameters in the `main.go` file.

  - num: number of nodes totally
  - psnum: the shard size of process shard
  - csnum: the shard size of control shard
  - adv: the runtime adversarial nodes ratio in different epoch. (If your server is not very powerful, I recommend you to use a smaller value, which can somewhat reduce the burden of concurrency on the server.)
  - T1: time bound 1
  - lamba: the lamba value used for calculate T2.

- run it by ./start.sh
- By default all nodes run locally, if you want to change the ip address of the node, you can do so in get drand.go to change the initialized address.
