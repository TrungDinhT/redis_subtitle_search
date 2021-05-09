OUT_REDIS_CMD = []

def read_file(filename):
    with open(filename, 'r') as f:
        for (i, line) in enumerate(f):
            line = line.strip('\n')
            if(i % 2 == 0):
                start_time = line
            else:
                OUT_REDIS_CMD.append("HSET " + "video_part:" + str(
                    len(OUT_REDIS_CMD)) + " start \"" + start_time + "\" subtitle \"" + line + "\"")

def write_file(filename):
    with open(filename, 'w') as f:
        for cmd in OUT_REDIS_CMD:
            f.write(cmd + '\n')

read_file('data/test1.txt')

print(OUT_REDIS_CMD)

write_file('redis-cmd/import_subs.redis')
