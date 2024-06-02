# JSON Parser / JQ

- https://codingchallenges.fyi/challenges/challenge-json-parser
- https://codingchallenges.fyi/challenges/challenge-jq/

```sh
# open the program in $(EDITOR)
make

# run the program on an input file
./file_parser.py test_json/full_suite/pass1.json

# run identify function (".") using jq
curl -s 'https://dummyjson.com/quotes?limit=2' | ./jq.py .

# run all tests
make test
# run specific test function
make test test=tests.test_states.TestStringState.test_basic_string
# run entire test class
make test test=tests.test_states.TestStringState
make test test=tests.test_json_struct.TestJSONStruct_pretty_print
# run entire test module
make test test=tests.test_states
```

TODO: add linter
