import csv
import argparse
import math
from pathlib import Path

parser = argparse.ArgumentParser()
parser.add_argument('-p', '--postclient', type=str, help='Output of the posting client.')
parser.add_argument('-g', '--getclient', type=str, help='Output of the getting client.')
parser.add_argument('-l', '--location', type=str, help='Location of the server.')
parser.add_argument('-f', '--file', type=str, help='File to save results.')
args = parser.parse_args()

def read_csv_file(file_path):
    column1_values = []
    column2_values = []
    with open(file_path, 'r') as file:
        reader = csv.reader(file)
        for row in reader:
            if len(row) >= 2:  
                column1_values.append(row[0])
                column2_values.append(row[1])
    return column1_values, column2_values

file1_path = args.postclient
file2_path = args.getclient

list1_col_bits, list1_col_time = read_csv_file(file1_path)

list2_col_bits, list2_col_time = read_csv_file(file2_path)

delays = []
time2_indexes = []
arrived_data = 0
last_bits = ""
bits1_index = 0
for index, bits2 in enumerate(list2_col_bits):
    if list1_col_bits[bits1_index] == bits2 or last_bits == bits2:
        time2_indexes.append(index)
        bits1_index += 1
        last_bits = ""
    else:
        last_bits = list1_col_bits[bits1_index]

time1_index = 0
for item in time2_indexes:
    delays.append((int(list2_col_time[item]) - int(list1_col_time[time1_index])) / 1000)
    time1_index += 1

print(delays)

if len(list1_col_time) != len(delays):
    print("Something is bad!")
    
delays.sort()
avg_delay = round((sum(delays) / len(delays)), 3)
delay_50th = delays[math.floor(len(delays) * 0.5)]
delay_95th = delays[math.floor(len(delays) * 0.95)]
delay_99th = delays[math.floor(len(delays) * 0.99)]

results = []
results.append(args.location)
results.append(avg_delay)
results.append(min(delays))
results.append(max(delays))
results.append(delay_50th)
results.append(delay_95th)
results.append(delay_99th)

if not Path("./" + args.file).exists():
    with open(args.file, mode='a', newline='') as file:
        file.write("location,avg_delay,min_delay,max_delay,median,95th,99th\n")
    
with open(args.file, mode='a', newline='') as file:
    writer = csv.writer(file)
    writer.writerow(results)
    
def sort_by_first_element(row):
    return row[0]

with open(args.file, mode='r') as file:
    reader = csv.reader(file)
    data = list(reader)

sorted_data = sorted(data[1:], key=sort_by_first_element)
sorted_data.insert(0, data[0])

with open(args.file, mode='w', newline='') as file:
    writer = csv.writer(file)
    writer.writerows(sorted_data)
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
 
