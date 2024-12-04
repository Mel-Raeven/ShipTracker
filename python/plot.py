import boto3
import folium
import pandas as pd
from folium import plugins
import datetime

# Step 1: Connect to DynamoDB
dynamodb_client = boto3.client('dynamodb', region_name='eu-central-1')  # Replace with your region

# Step 2: Scan DynamoDB for ship data
try:
    response = dynamodb_client.scan(TableName='ShipCords')  # Replace with your table name
    items = response['Items']

    # Step 3: Extract the coordinates and timestamps
    ships = []
    for item in items:
        ship_name = item['Name']['S']
        mmsi = item['MMSI']['N']
        latitude = float(item['Latitude']['N'])
        longitude = float(item['Longitude']['N'])
        timestamp = item['TS']['N']  # Assuming timestamp is a Unix epoch time (number of seconds)
        
        # Convert Unix timestamp to datetime
        timestamp = datetime.datetime.utcfromtimestamp(int(timestamp))

        ships.append({
            'ship_name': ship_name,
            'mmsi': mmsi,
            'latitude': latitude,
            'longitude': longitude,
            'timestamp': timestamp
        })

    # Convert list to DataFrame for easy handling
    ships_df = pd.DataFrame(ships)

    # Step 4: Sort the coordinates by timestamp
    ships_df = ships_df.sort_values(by='timestamp')

    # Step 5: Create a base map (you can set the initial location and zoom level)
    start_coords = [ships_df['latitude'].iloc[0], ships_df['longitude'].iloc[0]]
    folium_map = folium.Map(location=start_coords, zoom_start=5)

    # Step 6: Plot the tracks with lines and add timestamps as popups
    for ship_name, ship_data in ships_df.groupby('ship_name'):
        coordinates = list(zip(ship_data['latitude'], ship_data['longitude']))
        
        # Add a PolyLine for the ship's path
        folium.PolyLine(locations=coordinates, color='blue', weight=2.5, opacity=1).add_to(folium_map)

        # Add markers with timestamps as popups
        for idx, row in ship_data.iterrows():
            folium.Marker(
                location=[row['latitude'], row['longitude']],
                popup=f"Timestamp: {row['timestamp'].strftime('%Y-%m-%d %H:%M:%S')}",
                icon=folium.Icon(color='red', icon='info-sign')
            ).add_to(folium_map)

    # Step 7: Show the map
    folium_map.save('ship_tracks_map.html')
    print("Map has been saved as 'ship_tracks_map.html'.")

except Exception as e:
    print(f"Error occurred: {e}")
