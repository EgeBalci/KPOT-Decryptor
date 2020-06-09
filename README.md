# Kpot Decryptor
Tool for finding KPOT XOR key using known-plaintext attack. KPOT uses the XOR key to cipher requests between victim and C2. This script will recover the XOR key up to 27 bytes.

## Usage
- go build keypot.go
```
./keypot -u http://targetkpot.com/
```