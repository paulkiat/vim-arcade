# TCP/IP Suite
Folder/File          Description
---------------      ---------------
===============      ===============
---------------       ---------------
===============       ===============
    └── tests/
        ├── dhcp_test.go  Tests DHCP server operations, including IP assignments and lease management.
        ├── ip_routing_test.go  Tests static routing functionalities.
        ├── dynamic_routing_test.go  Tests dynamic routing protocol implementations.
        ├── multicast_test.go  Tests multicast group management and message handling.
        └── udp_service_test.go  Tests UDP Echo Server and Client interactions.

---------------      ------------------
===========           ===================
==============        ===================
    └── tests/         Tests DHCP server operations, including IP assignments and lease management.
    ├── **config.yaml**  Central configuration file for all services and modules, defining settings for DHCP, IP Routing, Dynamic Routing, Multicast, UDP Services, and logging preferences.
    ├── **dhcp_server/**
    │   └── **dhcp_server.go**: Implements the DHCP server using Go, handling DHCP Discover and Request messages, sending Offer and ACK responses, and managing IP leases through `lease_manager.go`.
    ├── **ip_routing/**
    │   └── **ip_routing.go**: Implements static routing logic, allowing the addition, removal, and display of static routes. Interfaces with the routing table to manage packet forwarding.
    ├── **dynamic_routing/**
    │   └── **dynamic_routing.go**: Implements dynamic routing protocols, handling route advertisements, learning from neighboring routers, and updating the routing table accordingly.
    ├── **multicast/**
    │   └── **multicast_service.go**: Implements multicast group management, providing an interactive command interface for sending messages, joining/leaving groups, and handling incoming multicast traffic.
    ├── **udp_service/**
    │   ├── **udp_service.go**: Implements the UDP Echo Server, echoing received messages back to the sender.
    │   └── **udp_client.go**: Implements a UDP Client for sending messages to the UDP Echo Server and receiving echoed responses, facilitating testing and interaction.
    ├── **lease_manager.go**: Handles DHCP lease assignments, renewals, expirations, and ensures efficient IP address utilization within the defined pool.
    ├── **packet_handler.go**: Contains utilities for parsing incoming network packets and constructing appropriate responses based on protocol specifications.
    └── **logger.go**: Sets up and manages logging across all modules, ensuring consistent and centralized logging for monitoring and troubleshooting.
    ├──  **dhcp_server/**
    │   └──  **dhcp_server.go**: Implements the DHCP server using Go, handling DHCP Discover and Request messages, sending Offer and ACK responses, and managing IP leases through `lease_manager.go`
    ├──  **ip_routing/**
    │   └──  **ip_routing.go**: Implements static routing logic, allowing the addition, removal, and display of static routes. Interfaces with the routing table to manage packet forwarding.
    ├──  **dynamic_routing/**
    │   └──  **dynamic_routing.go**: Implements dynamic routing protocols, handling route advertisements, learning from neighboring routers, and updating the routing table accordingly.
    ├──  **multicast/**
    │   └──  **multicast_service.go**: Implements multicast group management, providing an interactive command interface for sending messages, joining/leaving groups, and handling incoming multicast traffic.
    ├──  **udp_service/**
    │   ├──  **udp_service.go**: Implements the UDPEcho Server, echoing received messages back to the sender.
    │   └──  **udp_client.go**: Implements a UDP Client for sending messages to the UDPEcho Server and receiving echoed responses, facilitating testing and interaction.
    ├──  **lease_manager.go**: Handles DHCP lease assignments, renewals, expirations, and ensures efficient IP address utilization within the defined pool.
    ├──  **packet_handler.go**: Contains utilities for parsing incoming network packets and constructing appropriate responses based on protocol specifications.
    └──  **logger.go**: Sets up and manages logging across all modules, ensuring consistent and centralized logging for monitoring and troubleshooting

# Detailed Folder and File Descriptions
1. /tcpip_suite/
1. `/tcpip_suite/`

3. dhcp_server/
1.  /tcpip_suite/
    1.  `tcpip_suite/`
2. `config.yaml`
   • Description: Centralized configuration file in YAML format, defining settings for DHCP, IP Routing, Dynamic Routing, Multicast, UDP Services, and logging preferences.
3. `dhcp_server/`
   • Description: Contains all components related to the DHCP server, responsible for dynamic IP address assignment, lease management, and DHCP option provisioning.
   • dhcp_server.go:
     • Functionality: Implements the DHCP server using Go, handling DHCP Discover and Request messages, sending Offer and ACK responses, and managing IP leases through lease_manager.go.

4. ip_routing/
4. `ip_routing/`

	•	Description: Contains utilities for parsing incoming network packets and constructing appropriate responses based on protocol specifications.
	•	Functionality: Facilitates the interpretation of DHCP, multicast, and UDP packets, enabling the suite to respond correctly to network events.

10. logger.go
        • Functionality: Implements the DHCP server using Go, handling DHCP Discover and Request messages, sending Offer and ACK responses, and managing IP leases through `lease_manager.go`.
4.  ip_routing/
    1. `ip_routing/`
        2.	Description: Contains utilities for parsing incoming network packets and constructing appropriate responses based on protocol specifications.
        3.	Functionality: Facilitates the interpretation of DHCP, multicast, and UDP packets, enabling the suite to respond correctly to network events.
5. logger.go

	•	Description: Sets up and manages logging across all modules, ensuring consistent and centralized logging for monitoring and troubleshooting.
	•	Functionality: Configures logging levels, formats, and output destinations based on config.yaml, and provides logging utilities for other modules.

11. tests/

	•	Description: Directory containing unit and integration tests for all networking modules, ensuring reliability and correctness of functionalities.
	•	dhcp_test.go: Tests DHCP server operations, including IP assignments and lease management.
	•	ip_routing_test.go: Tests static routing functionalities.
	•	dynamic_routing_test.go: Tests dynamic routing protocol implementations.
	•	multicast_test.go: Tests multicast group management and message handling.
	•	udp_service_test.go: Tests UDP Echo Server and Client interactions.

12. README.md

	•	Description: Provides comprehensive documentation, including an overview of the suite, setup instructions, usage guidelines, feature descriptions, and contribution protocols.

13. LICENSE

	•	Description: Details the licensing terms under which the TCP/IP suite is distributed, specifying usage rights and restrictions.

14. scripts/

	•	Description: Contains utility scripts for automating deployment, backups, and monitoring of the TCP/IP suite.
	•	deploy.sh: Automates the deployment process, setting up necessary services and configurations.
	•	backup.sh: Automates the backup of configurations, logs, and lease data.
	•	monitor.sh: Integrates with monitoring tools to oversee the health and performance of the suite.
5. `dynamic_routing/`
    4. Description: Sets up and manages logging across all modules, ensuring consistent and centralized logging for monitoring and troubleshooting.
    5. Functionality: Configures logging levels, formats, and output destinations based on config.yaml, and provides logging utilities for other modules.
6.  `dynamic_routing/`
    ├── **dhcp_server/**  ├─ Implements the DHCP server using Go, handling DHCP Discover and Request messages, sending Offer and ACK responses, and managing IP leases through `lease_manager.go`.
    ├── **ip_routing/**    │ └── **ip_routing.go**: Implements static routing logic, allowing the addition, removal, and display of static routes. Interfaces with the routing table to manage packet forwarding.
    ├── **dynamic_routing/**  ├── **dynamic_routing.go**: Implements dynamic routing protocols, handling route advertisements, learning from neighboring routers, and updating the routing table accordingly.
    ├── **multicast/**      │   └── **multicast_service.go**: Implements multicast group management, providing an interactive command interface for sending messages, joining/leaving groups, and handling incoming multicast traffic.
    ├── **udp_service/**      ├── **udp_service.go**: Implements the UDPEcho Server, echoing received messages back to the sender.
                      │   └── **udp_client.go**: Implements a UDP Client for sending messages to the UDPEcho Server and receiving echoed responses, facilitating testing and interaction.
    ├── **lease_manager.go**      ── Handles DHCP lease assignments, renewals, expirations, and ensures efficient IP address utilization within the defined pool.
    ├── **packet_handler.go**       ── Contains utilities for parsing incoming network packets and constructing appropriate responses based on protocol specifications.
    └── **logger.go**                ────── Sets up and manages logging across all modules, ensuring consistent and centralized logging for monitoring and troubleshooting.
   • Description: Facilitates dynamic routing protocols (e.g., RIP, OSPF) to enable automatic route discovery and adaptation to network topology changes.
6. `multicast/`
   • Description: Manages multicast group subscriptions, allowing devices to join or leave multicast groups and handle the sending and receiving of multicast messages.
7. `udp_service/`
   • Description: Provides UDP-based services, including an Echo Server for testing and communication purposes.
8. `lease_manager.go`
   • Description: Handles DHCP lease assignments, renewals, expirations, and ensures efficient IP address utilization within the defined pool.
9. `packet_handler.go`
7.  `multicast/`
    1.	Description: Manages multicast group subscriptions, allowing devices to join or leave multicast groups and handle the sending and receiving of multicast messages.
8.  udp_service/
    • Description: Provides UDP-based services, including anEcho Server for testing and communication purposes.
9. `lease_manager.go`
    10.	Description: Handles DHCP lease assignments, renewals, expirations, and ensures efficient IP address utilization within the defined pool.
10.  packet_handler.go
   • Description: Contains utilities for parsing incoming network packets and constructing appropriate responses based on protocol specifications.
10. `logger.go`
    • Description: Sets up and manages logging across all modules, ensuring consistent and centralized logging for monitoring and troubleshooting.

# Tests
11. `tests/`
    • Description: Directory containing unit and integration tests for all networking modules, ensuring reliability and correctness of functionalities.
12. `README.md`
    • Description: Provides comprehensive documentation, including an overview of the suite, setup instructions, usage guidelines, feature descriptions, and contribution protocols.

# Licensing
13. `LICENSE`
    • Description: Details the licensing terms under which the TCP/IP suite is distributed, specifying usage rights and restrictions.

# Scripts
14. `scripts/`
    • Description: Contains utility scripts for automating deployment, backups, and monitoring of the TCP/IP suite.
15. `deploy.sh`: Automates the deployment process, setting up necessary services and configurations.
16. `backup.sh`: Automates the backup of configurations, logs, and lease data.
17. `monitor.sh`: Integrates with monitoring tools to oversee the health and performance of the suite.

5. `lease_manager.go`
    12.	Description: Handles DHCP lease assignments, renewals, expirations, and ensures efficient IP address utilization within the defined pool.
13.  packet_handler.go
    • Description
