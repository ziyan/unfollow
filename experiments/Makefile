.PHONY: all
all: data experiments

CC = g++
CFLAGS = -Isnap/snap-core -Isnap/glib-core -Wall -O4

experiments: experiments.o snap/snap-core/Snap.o
	$(CC) -o $@ $^

%.o: %.cpp
	$(CC) $(CFLAGS) -c -o $@ $<

snap/snap-core/Snap.o: snap
	cd snap && make all

.PHONY: data
data:
	cd data && make

.PHONY: clean
clean:
	rm -f experiments
	rm -f *.o

