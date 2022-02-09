# Trigrams

Generates random text based on trigrams generated from input text

## Contents

- [Building](#building)
- [Running](#running)
- [Using](#using)
- [Implementation notes](#implementation-notes)
  * [NGram size](#ngram-size)
  * [Maximum word count](#maximum-word-count)
  * [Punctuation stripping](#punctuation-stripping)
  * [Weighted random selection](#weighted-random-selection)
  * [Endpoint considerations](#endpoint-considerations)
  * [ioutil.ReadAll() vs streaming requests](#ioutilreadall-vs-streaming-requests)
    + [ioutil.ReadAll()](#ioutilreadall)
    + [Streaming](#streaming)
- [Testing](#testing)
  * [Running unit tests](#running-unit-tests)
  * [Race detection](#race-detection)
  * [Buffer read size vs learn request time](#buffer-read-size-vs-learn-request-time)

## Building

```go build```

## Running

```./trigrams```

## Using

Train the server with some text by doing:

```curl -X POST -d "this is some text" http://localhost:8080/learn```

Text files can be sent as follows (make sure to use --data-binary to account for new lines):

```curl -X POST --data-binary @pride-prejudice.txt http://localhost:8080/learn```

Generate a random string of text by running:

```curl -X GET http://localhost:8080/generate```

Output:

``` claimed towards Mr. Darcy had never seen a collection of people in this manner; and as a rector, made him altogether a mixture of pride and impertinence; she had as good a chance of happiness as if the second, I can admire you much better finish his letter. When that business was over, he applied to Miss Grantley’s.” “Will you give me leave to apologise for it, as well as her mother should be in danger of hating each other for the other. The master of the impertinent. She mentioned this to her notice. Mrs. Phillips was quite disconcerted. She ```

## Implementation notes

The application is controlled by a series of consts, though these could easily be replaced with command line flags.

### NGram size

Although the challenge focuses on trigrams, there are actually a number of different ngram sizes that can be processed;
unigrams. bigrams, trigrams etc

The application can be developed with a more general use case; specifically, making the gram size variable. Hence, we'll
develop an ngram parsing application, configured to parse trigrams.

### Maximum word count

The possibility of generating unigrams highlights the need for a maximum word count; we don't want to fall into an
infinite loop when building random strings. Additionally, it cannot be ruled out that we won't fall into an infinite
loop when parsing larger grams, so maximum word count will be controlled by a variable.

### Punctuation stripping

Input text will typically be composed of alphanumeric characters, and any number of punctuation characters. I'm unsure
if this task wants to include or omit such punctuation from the ngrams; hence, stripping punctuation from strings will
be controlled by a variable also.

### Weighted random selection

Because we want the random word selection to be reflective of the frequency with which that word occurs in the source
text, we need an efficient way of weighting the random selection. A naive approach would be to create a new array of
strings, duplicating each entry a number of times equal to that word's frequency, but this would be horribly inefficient.
Instead, if we calculate the total frequency count (F) of all words in the source text, then generate a random number
(R) between 1 and (F), and subtract each word's frequency from (R) at random, then test for the current value of (R)
being less than or equal to 0, we'll have selected a random word based on that word's frequency. For example, given:

- "dog"  : 5
- "cat"  : 1
- "bird" : 3
- "snake": 1
- "horse": 2

The total frequency is 5 + 1 + 3 + 1 + 2 = 12. Generating a random number between 1 and 12 might give us R = 7. Randomly
deleting each frequency from (R) until (R) is <= 0 might give us:

7 - (cat : 1) = 6
<= 0 : NO
6 - (bird : 3) = 3
<= 0 : NO
3 - (dog : 5) = -2
<= 0 : YES
dog is selected

Or, another example where R = 5:

5 - (dog : 5) = 0
<= 0 : YES
dog is selected

### Endpoint considerations

A naive approach with endpoints is to call a time-consuming function asynchronously and respond immediately with an OK
status code.

The problem with this approach is that flooding this endpoint with a large number of calls will cause the application to
crash.

Instead of unbounded async go function calls which will quickly exhaust resources, we need to use channels to limit
access to resources

We have two queues:
  - learn queue
  - generation queue

Both queues accept tasks; the learn queue accepts learn tasks, and the generation queue accepts generation tasks

Next, we launch dispatchers

The learn dispatcher creates a pool of learn workers, and starts them. Each worker runs, and registers itself with the
dispatch pool, then waits to receive a learning task.

The dispatcher then listens to the learn queue, waiting for new learn tasks. When a new task is received, a worker is
retrieved from the worker pool, and the task is passed to that worker to be processed.

Once a task is received by the worker, the task is processed and the input text is parsed into tokens, then n-grams are
gathered and added to the gram collection.

The same happens for the generation queue for generate tasks.

### ioutil.ReadAll() vs streaming requests

Initially, the `/learn` endpoint read the entire request body into memory for processing. While this worked well for small requests, copying very large requests into memory at once might cause the application to crash. An alternative solution is to stream the request body, rather than copy the entire request body into memory at once. To test the new implementation, five concurrent requests were made to the `/learn` endpoint, with [enwiki8](http://mattmahoney.net/dc/textdata.html) (approx. 95MB) as the request body. In the case of the `ioutil.ReadAll()` implementation, memory usage of the application almost immediately spiked to over 3GB. In the streaming implementation, under the same use case, memory usage climbed only to ~30MB in seconds, and ~60MB in minutes:

#### ioutil.ReadAll()

![ReadAll](/readall.png)

#### Streaming

![Streaming](/streaming.png)  

## Testing

### Running unit tests

Run unit tests with:

```go test ./...```

### Race detection

Build the program with the `-race` flag:

```go build -race```

Running the application:

```./trigrams```

And then passing some data to the `/learn` endpoint:

```curl -X POST --data-binary @pride-prejudice.txt http://localhost:8080/learn```

Followed by a few hundred `/generate` calls:

```hey -n 200 -c 10 http://localhost:8080/generate```

Does not report any race conditions.

### Buffer read size vs learn request time

In order to identify the optimum size (in bytes) of a buffer for reading the request body, `pride-prejudice.txt` was submitted to the `/learn` endpoint, each time with a different `ReadSize` specified as the length of the buffer for reading the request body. Over the following series of tests, a buffer size of 64 bytes was identified as achieving the best request read performance:

| read size (bytes) | request time |
|-----------|--------------|
|   1024    |   1m2.536s   |
|   1       |   1m16.315s  |
|   2048    |   1m7.439s   |
|   512     |   0m57.830s  |
|   256     |   0m55.934s  |
|   128     |   0m54.770s  |
|  **64**     |  **0m54.367s** |
|   32      |   0m56.060s  |
|   48      |   0m55.345s  |
|   56      |   0m54.730s  |
|   60      |   0m54.509s  |
|   62      |   0m54.741s  |
|   63      |   0m54.724s  |
