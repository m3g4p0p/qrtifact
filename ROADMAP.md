# Roadmap

- Simple local P2P connection using listener / dialer model (e.g. Unix socket) with state sync (e.g. shared counter)
    - Consider JSON-RPC 2.0 as wire protocol
- Implement file-based transport to replace socket
    - Watch inbox directory for new files
    - Read new files as incoming messages
    - Write messages to outbox directory
    - Share outbox with peer
- Implement QR codec for file transport
- Implement handshake with key exchange for encryption / signed messages
- Implement actual game logic