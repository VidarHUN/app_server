import csv
import argparse
import math
import pandas as pd
import tqdm
from pathlib import Path

parser = argparse.ArgumentParser()
parser.add_argument('-p', '--postclient', type=str, help='Output of the posting client.')
parser.add_argument('-g', '--getclient', type=str, help='Output of the getting client.')
parser.add_argument('-l', '--location', type=str, help='Location of the server.')
parser.add_argument('-f', '--file', type=str, help='File to save results.')
args = parser.parse_args()

def row_compare(row1, row2, idx1, idx2):
    try:
        # print(row1)
        # print(idx1)
        # print(row2)
        # print(idx2)
        if (row1.at[idx1, 'length'] == row2.at[idx2, 'length'] and
            row1.at[idx1, 'group_id'] == row2.at[idx2, 'group_id'] and
            row1.at[idx1, 'object_id'] == row2.at[idx2, 'object_id']):
            return True

        return False

    except Exception as e:
        # Handle the exception here.
        # For example, you could log the exception or print a message to the user.
        print(row1)
        print(idx1)
        print(row2)
        print(idx2)
        pass


def sort_csv_files(file1, file2):
    """Sorts two CSV files by group ID and object ID without column names.

    Args:
        file1: Path to the first CSV file.
        file2: Path to the second CSV file.
    """

    # Read the CSV files into Pandas DataFrames using the `header=None` parameter.
    df1 = pd.read_csv(file1, header=None, delimiter=',')
    df2 = pd.read_csv(file2, header=None, delimiter=',')

    # Rename the columns to `group_id`, `object_id`, `last_24_byte_of_sent_message`, `unix_timestamp_in_us`, and `packet_length`.
    df1.columns = ['bytes', 'ts', 'length', 'group_id', 'object_id']
    df2.columns = ['bytes', 'ts', 'length', 'group_id', 'object_id']

    df1 = df1.drop('bytes', axis=1)
    df2 = df2.drop('bytes', axis=1)

    # Sort the DataFrames by group ID and object ID.
    df1 = df1.sort_values(by=['group_id', 'object_id'], ascending=True, ignore_index=True)
    df2 = df2.sort_values(by=['group_id', 'object_id'], ascending=True, ignore_index=True)

    # Write the sorted DataFrames back to CSV files.
    df1.to_csv('post.csv', index=False, header=False)
    df2.to_csv('get.csv', index=False, header=False)

    # Create two new DataFrames to store the filtered rows.
    new_df1 = pd.DataFrame(columns=['ts', 'length', 'group_id', 'object_id'])
    new_df2 = pd.DataFrame(columns=['ts', 'length', 'group_id', 'object_id'])

    lost_packets = 0
    bigger = 0
    smaller = 0
    first = True

    index_df1 = 0
    index_df2 = 0

    max_group_id = df1['group_id'].max()
    groups1 = []
    groups2 = []
    for _ in range(max_group_id + 1):
        groups1.append([])
        groups2.append([])

    for _, row in df1.iterrows():
        group_id = row["group_id"]
        groups1[group_id].append(row)

    for _, row in df2.iterrows():
        group_id = row["group_id"]
        groups2[group_id].append(row)

    # pbar = tqdm.tqdm(total=len(groups1), desc="Process csvs")
    for idx, group1 in enumerate(tqdm.tqdm(groups1)):
        # print(idx)
        group2 = groups2[idx]
        index_df1 = 0
        index_df2 = 0
        next_row = None
        while index_df1 < len(group1) and index_df2 < len(group2):
            row1 = group1[index_df1]
            row2 = group2[index_df2]
            next_row = row2
            if row1['length'] == row2['length']:
                new_df1 = pd.concat([new_df1, row1.to_frame().T], ignore_index=True)
                new_df2 = pd.concat([new_df2, row2.to_frame().T], ignore_index=True)
                index_df1 = index_df1 + 1
                index_df2 = index_df2 + 1
                # pbar.update(1)
            else:
                nr_rows = 1
                ts = row2['ts']
                length = row2['length']

                while length < row1['length'] and index_df2 < len(group2) - 1 and nr_rows < 3:
                    index_df2 += 1
                    next_row = group2[index_df2]
                    nr_rows += 1
                    ts += next_row['ts']
                    length += next_row['length']

                ts = math.floor(ts / nr_rows)

                if length == row1['length']:
                    new_row = pd.Series({'ts': ts, 'length': length, 'group_id': row2['group_id'], 'object_id': row2['object_id']})
                    new_df1 = pd.concat([new_df1, row1.to_frame().T], ignore_index=True)
                    new_df2 = pd.concat([new_df2, new_row.to_frame().T], ignore_index=True)
                    index_df1 = index_df1 + 1
                    index_df2 = index_df2 + 1
                    # pbar.update(1)
                else:
                    # Drop the whole object if the arriving order messed up
                    object_id = row1['object_id']
                    while object_id == row1['object_id'] and index_df1 < len(group1) - 1:
                        index_df1 += 1
                        row1 = group1[index_df1]

                    row2 = next_row

                    while object_id == row2['object_id'] and index_df2 < len(group2) - 1:
                        index_df2 += 1
                        row2 = group2[index_df1]

                    # if length > row1['length']:
                    #     bigger += 1
                    # else:
                    #     smaller += 1
                    # if first:
                    #     print(row2.to_string())
                    #     # print(next_row.to_string())
                    #     print(row1.to_string())
                    #     first = False
                    index_df1 = index_df1 + 1
                    index_df2 = index_df2 + 1
                    # pbar.update(1)


    # while index_df1 < len(df1) or index_df2 < len(df2):
    #     # print(index_df1, index_df2)
    #     row1 = df1.iloc[[index_df1]]
    #     row2 = df2.iloc[[index_df2]]
    #     if row_compare(row1, row2, index_df1, index_df2):
    #         new_df1 = pd.concat([new_df1, row1], ignore_index=True)
    #         new_df2 = pd.concat([new_df2, row2], ignore_index=True)
    #         index_df1 = index_df1 + 1
    #         index_df2 = index_df2 + 1
    #         pbar.update(1)
    #     else:
    #         nr_rows = 1
    #         ts = row2.at[index_df2, 'ts']
    #         length = row2.at[index_df2, 'length']
    #         group_id = row2.at[index_df2, 'group_id']
    #         object_id = row2.at[index_df2, 'object_id']

    #         # next_row = None
    #         # if length > row1.at[index_df1, 'length']:
    #         #     lost_packets = lost_packets + 1
    #         #     index_df1 = index_df1 + 1
    #         #     pbar.update(1)
    #         #     continue
    #         # while length < row1.at[index_df1, 'length']:
    #             # print(nr_rows, length, row1.at[index_df1, 'length'])
    #         index_df2 = index_df2 + 1
    #         if index_df2 >= len(df2):
    #             break
    #         next_row = df2.iloc[[index_df2]]
    #         # print(next_row)
    #         length = length + next_row.at[index_df2, 'length']
    #         ts = ts + next_row.at[index_df2, 'ts']
    #         nr_rows = nr_rows + 1

    #         ts = math.floor(ts / nr_rows)
    #         if length == row1.at[index_df1, 'length']:
    #             new_row = {'ts': ts, 'length': length, 'group_id': group_id, 'object_id': object_id}
    #             tmp_df = pd.DataFrame(columns=['ts', 'length', 'group_id', 'object_id'])
    #             tmp_df.loc[0] = new_row
    #             if row_compare(row1, tmp_df, index_df1, 0):
    #                 new_df1 = pd.concat([new_df1, row1], ignore_index=True)
    #                 new_df2 = pd.concat([new_df2, tmp_df], ignore_index=True)
    #                 index_df1 = index_df1 + 1
    #                 index_df2 = index_df2 + 1
    #                 pbar.update(1)
    #             else:
    #                 index_df1 = index_df1 + 1
    #                 pbar.update(1)
    #         else:
    #             if length > row1.at[index_df1, 'length']:
    #                 bigger += 1
    #             else:
    #                 smaller += 1
    #             if first:
    #                 print(row2.to_string())
    #                 print(next_row.to_string())
    #                 print(row1.to_string())
    #                 first = False
    #             index_df1 = index_df1 + 1
    #             index_df2 = index_df2 + 1
    #             pbar.update(1)
    #         # elif group_id != next_row.at[index_df2, 'group_id'] or object_id != next_row.at[index_df2, 'object_id']:
    #         #     # print('not_matching')
    #         #     lost_packets = lost_packets + 1
    #         #     # index_df1 = index_df1 + 1
    #         #     # pbar.update(1)
    #         # else:
    #         #     print(row1.to_string())
    #         #     print(next_row.to_string())
    #         #     index_df1 = index_df1 + 1
    #         #     index_df2 = index_df2 + 1
    #         #     pbar.update(1)


    print(len(new_df1.index))
    # print(f'bigger: {bigger}')
    # print(f'smaller: {smaller}')

    # Write the sorted DataFrames back to CSV files.
    new_df1.to_csv('post.csv', index=False, header=False)
    new_df2.to_csv('get.csv', index=False, header=False)

def read_csv_file(file_path):
    column1_values = []
    column2_values = []
    with open(file_path, 'r') as file:
        reader = csv.reader(file)
        pbar = tqdm.tqdm(total=reader.line_num, desc="Read csvs")
        for row in reader:
            pbar.update(1)
            if len(row) >= 2:
                column1_values.append(row[0])
    return column1_values

file1_path = args.postclient
file2_path = args.getclient

sort_csv_files(file1_path, file2_path)

list1_col_time = read_csv_file('post.csv')

list2_col_time = read_csv_file('get.csv')

delays = []
time2_indexes = []
arrived_data = 0
last_bits = ""
bits1_index = 0
# for index, bits2 in enumerate(list2_col_bits):
#     if list1_col_bits[bits1_index] == bits2 or last_bits == bits2:
#         time2_indexes.append(index)
#         bits1_index += 1
#         last_bits = ""
#     else:
#         last_bits = list1_col_bits[bits1_index]

pbar_delays = tqdm.tqdm(total=len(list1_col_time), desc="Calculate delays")
for idx, _ in enumerate(list1_col_time):
    pbar_delays.update(1)
    delay = (int(list2_col_time[idx]) - int(list1_col_time[idx])) / 1000
    if delay > 0:
        delays.append(delay)

# print(delays)

# if len(list1_col_time) != len(delays):
#     print(len(list1_col_time))
#     print(len(delays))
#     print("Something is bad!")

pbar_diffs = tqdm.tqdm(total=len(delays), desc="Calculate jitter")
differences = []
for i in range(len(delays) - 1):
    pbar_diffs.update(1)
    difference = abs(delays[i + 1] - delays[i])
    differences.append(difference)

# Calculate the average jitter.
jitter = sum(differences) / len(differences)

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
results.append(jitter)

if not Path("./" + args.file).exists():
    with open(args.file, mode='a', newline='') as file:
        file.write("location,avg_delay,min_delay,max_delay,median,95th,99th,jitter\n")

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