## Install
=======

    - Make sure you have the AWS credentials set either in ~/.aws/credentials or in a bash export
    - Put the sqsutils executable in your path somewhere...
    - Make sure it is executable (chmod +x sqsutils)
    - Done

    run help:  `sqsutils -h`


## Example uses:
=============

 ./sqsutils -a depth -s sqs.us-east-1.amazonaws.com -i 478989820108 DEV_TIBCO_STORES_QUEUE
 	-s set to amazon server + amazonId set will direct requests to AWS/SQS (for real - becareful)
 	(depth is the default so action is not required)
 	-i amazon ID for amazon requests
 	-a actions
 		-- depth (default)
 		-- create (make a queue)
 		-- purge (remove all messages from a queue)
 		-- send (put a test message on the queue)
 		-- receive (gets a message from the queue)
 	-t string body (argument for send - will send the string in this parameter in the message body
 	-f file body (argument for send - will send the contents of the file parameter in the message body

 ./sqsutils -a depth -s sqs.us-east-1.amazonaws.com -i 478989820108 DEV_TIBCO_STORES_QUEUE   // Talk to a real amazon queue

 ./sqsutils test_queue_1   // Gets the queue depth from local elastic mq queue test_queue_1

