import os
import mariadb
import sys
import uuid
import csv
from dotenv import load_dotenv
load_dotenv()
user = os.getenv('DB_USER')
passw = os.getenv('DB_PASSWORD')
db_name = os.getenv("DB_NAME")

try:
    conn = mariadb.connect(
        user=user,
        password=passw,
        host="127.0.0.1",
        port=3306,
        database=db_name
    )

    data = []
    with open("data.csv", "r") as csv_file:
        reader = csv.reader(csv_file)

        for i, v in enumerate(reader):
            obj = {}
            if i == 0:
                continue
            else:
                obj["VendorID"] = v[0]
                obj["VendorName"] = v[1]
                obj["Poster"] = v[2]
                obj["Date"] = v[3]
                obj["Rating"] = v[4]
                obj["Source"] = v[5]
                obj["Content"] = v[6]
                obj["LinkToSource"] = v[7]
            data.append(obj)

    cursor = conn.cursor()
    
    for obj in data:
        cursor.execute(
            "INSERT INTO vendor_review(review_id, vendor_id, vendor_name, poster, date, rating, source, content, link_to_source) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
            (str(uuid.uuid4()),
                obj["VendorID"],
                obj["VendorName"],
                obj["Poster"],
                obj["Date"],
                obj["Rating"],
                obj["Source"],
                obj["Content"],
                obj["LinkToSource"])
        )
        conn.commit()   
    print("Finish")
    conn.close()

except mariadb.Error as e:
    print(f"Error connecting to MariaDB Platform: {e}")
    sys.exit(1)
