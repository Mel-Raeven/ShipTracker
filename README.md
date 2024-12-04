# ShipTracker

ShipTracker is a project designed to monitor and track fishing boats to detect deviations from their normal routes, potentially identifying illicit activities such as drug pickups. The project leverages Go for tracking and storing data, AWS DynamoDB for database storage, and Python for visualizing the collected coordinates.

# About the Project

ShipTracker was developed as part of my Cyberstars Minor program. Its primary goal is to assist in identifying suspicious maritime activity by tracking fishing vessels. By analyzing deviations in their usual routes, ShipTracker can help flag potential illegal actions, such as the transportation of drugs or contraband.
Features

- Real-Time Tracking: Retrieves and tracks fishing boats' coordinates using Go.
- Data Storage: Efficiently stores location data in AWS DynamoDB for easy   querying and scalability.
- Route Visualization: Plots the tracked coordinates using Python, providing c- lear visual representation.
- Anomaly Detection: Flags significant route deviations for further inspection.

# Technologies Used

- Programming Languages:
    - Go – For data retrieval and backend processing.
    - Python – For data visualization and plotting.
- Database:
    - AWS DynamoDB – To store and manage the coordinates.
