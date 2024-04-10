# this is a python script to convert the given csv files to a sqlite3 database

import pandas as pd
from sqlalchemy import create_engine
from datetime import datetime, timezone
import pytz

start = datetime.now()

engine = create_engine("sqlite:///data.db")

df1 = pd.read_csv("1.csv")
df2 = pd.read_csv("2.csv")
df3 = pd.read_csv("3.csv")

def f1(x):
    x = x[:-4]
    tformat = "%Y-%m-%d %H:%M:%S.%f"
    try:
        x = datetime.strptime(x, tformat)
    except ValueError:  
        tformat = "%Y-%m-%d %H:%M:%S"
        x = datetime.strp(x, tformat)
    finally:
        return x

poll_time = df1["timestamp_utc"]

def f2(x):
    if x == "active": return 1
    else: return 0
    
df1['store_id'] = pd.to_numeric(df1['store_id'])
df1["status"] = df1["status"].apply(f2)
df1["timestamp_utc"] = poll_time.apply(f1)
df1 = df1.sort_values(['store_id', 'timestamp_utc'], ascending=[1, 0])
df1['store_id'] = df1['store_id'].apply(str)

df2 = df2.sort_values(['store_id', 'day', 'start_time_local', 'end_time_local'])

df1.to_sql("polls", con=engine, if_exists='replace', index=False)
df2.sort_values(by=['store_id', 'day'])
df2.to_sql("store_data", con=engine, if_exists='replace', index=False)
df3.to_sql("timezone", con=engine, if_exists='replace', index=False)
print(datetime.now()-start)