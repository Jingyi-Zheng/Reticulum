import socket
import time
import threading

# Set parameters
packet_size = 2048  # Size of each packet (bytes)
port_start = 10000  # Starting port number
# Modify this number
port_end = port_start + 249  # Ending port number
file_size = 2 * 1024 * 1024  # Size of the file to send (bytes)
num_threads = 1  # Number of threads for sending data

# Generate random data
data = bytearray(packet_size)
for i in range(packet_size):
    data[i] = i % 256

# Create a listening thread


def listen_thread(port):
    # Create a socket and start listening on the port
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.bind(('localhost', port))
    s.listen(1)
    conn, addr = s.accept()

    # Receive data and measure time
    start_time = time.time()
    recv_size = 0
    while recv_size < file_size:
        data_recv = conn.recv(packet_size)
        if not data_recv:
            break
        recv_size += len(data_recv)
    end_time = time.time()

    # Output the results
    elapsed_time = end_time - start_time
    print(f'Received total data size on port {port}: {recv_size} bytes, time: {elapsed_time:.2f}s')

    # Close the socket
    conn.close()

# Thread function for sending data


def send_thread(start_port):
    for port in range(start_port, port_end+1, num_threads):
        # Create a socket and connect to the port
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.connect(('localhost', port))
        print("Connected")
        # Send data
        send_size = 0
        while send_size < file_size:
            s.send(data)
            send_size += packet_size

        # Close the socket
        s.close()


# Start the listening threads
for port in range(port_start, port_end+1):
    t = threading.Thread(target=listen_thread, args=(port,))
    t.start()

# Start the threads for sending data and measure time
start_time = time.time()
send_threads = []
for i in range(num_threads):
    t = threading.Thread(target=send_thread, args=(port_start+i,))
    t.start()
    send_threads.append(t)

# Wait for all send data threads to complete and measure time
for t in send_threads:
    t.join()
end_time = time.time()

# Output the results
elapsed_time = end_time - start_time
print(f'Total time elapsed: {elapsed_time:.2f}s, average speed: {(file_size*num_threads*8/1024/1024)/elapsed_time:.2f} Mbps')
print('Data transmission completed')
