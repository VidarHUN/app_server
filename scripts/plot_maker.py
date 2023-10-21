import argparse
import csv
import statistics
import matplotlib.pyplot as plt
from matplotlib.colors import ListedColormap
import pandas as pd
import math

def final_avg(filename):
    newrow, delays, min_delays, max_delays, medians, delay_95ths, delay_99ths = [],[],[],[],[],[],[]
    location = None
    
    with open("final_" + filename, mode='a', newline='') as file:
        file.write("location,avg_delay,min_delay,max_delay,median,95th,99th\n")
    	
    with open(filename, 'r') as file:
        reader = csv.reader(file)
        next(reader)
        for row in reader:
        
            if row[0] != location and location is not None:
                with open("final_" + filename, mode='a', newline='') as file2:
                    writer = csv.writer(file2)
                    newrow.append([location, round(statistics.mean(delays), 3), min(min_delays), max(max_delays), round(statistics.mean(medians), 3), round(statistics.mean(delay_95ths), 3), round(statistics.mean(delay_99ths), 3)])
                    writer.writerows(newrow)
                newrow, delays, min_delays, max_delays, medians, delay_95ths, delay_99ths = [],[],[],[],[],[],[]

            location = row[0]
            delays.append(float(row[1]))
            min_delays.append(float(row[2]))
            max_delays.append(float(row[3]))
            medians.append(float(row[4]))
            delay_95ths.append(float(row[5]))
            delay_99ths.append(float(row[6]))
            
        with open("final_" + filename, mode='a', newline='') as file2:
           writer = csv.writer(file2)
           newrow.append([location, statistics.mean(delays), min(min_delays), max(max_delays), statistics.mean(medians), statistics.mean(delay_95ths), statistics.mean(delay_99ths)])
           writer.writerows(newrow)
           
    return "final_" + filename
                

parser = argparse.ArgumentParser()
parser.add_argument('-f', '--file', type=str, help='Data csv file.')
parser.add_argument('-n', '--name', type=str, help='Name of the plot.')
args = parser.parse_args()

cmap = ListedColormap(['#0343df', '#e50000', '#000000', '#929591', '#008000', '#964b00'])

final_file_name = final_avg(args.file)
df = pd.read_csv(final_file_name)

ax = df.plot.bar(x='location', colormap=cmap)

ax.set_xlabel('Location')
ax.set_ylabel('ms', rotation = 0)
ax.set_title(args.name)
plt.xticks(rotation=0)

resolution = 300  # Adjust the resolution as needed
plt.savefig(args.name + '.png', dpi=resolution)
