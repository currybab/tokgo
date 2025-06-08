import os
import tiktoken
import csv

encoding_types = ["cl100k_base", "p50k_base", "p50k_edit", "r50k_base", "o200k_base"]

for encoding_type in encoding_types:
    encoding = tiktoken.get_encoding(encoding_type)

    with open("../resources/test/base_prompts.csv", mode="r", encoding="utf-8") as f:
        csvdata = csv.reader(f, delimiter=",", quotechar='"')
        next(csvdata, None)
        with open(f"../resources/test/{encoding_type}_encodings.csv", mode="w", encoding="utf-8") as outFile:
            writer = csv.writer(outFile, delimiter=",", quotechar='"')
            writer.writerow(["input","output","outputMaxTokens10"])
            for row in csvdata:
                encoded = encoding.encode_ordinary(row[0])
                for i in reversed(range(11)):
                    encodedShort = encoded[:i]
                    decoded = encoding.decode(encodedShort)
                    if (row[0].startswith(decoded)):
                        writer.writerow([row[0], encoded, encodedShort])
                        break
