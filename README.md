## Install

    - Make sure you have the AWS credentials set either in ~/.aws/credentials or in a bash export
    - Put the sqsutils executable in your path somewhere...
    - Make sure it is executable (chmod +x sqsutils)
    - Done

    run help:  `sqsutils -h`


## Example uses:

    1.  ./sqsutils -a depth -s sqs.us-east-1.amazonaws.com -i 478989820108 STORES_QUEUE
 	
 	    -s set to amazon server + amazonId set will direct requests to AWS/SQS (for real - be careful)
 	    -i amazon ID for amazon requests
 	    -a actions
            -- depth (default so action if action is omitted)
            -- create (make a queue)
            -- purge (remove all messages from a queue)
            -- send (put a test message on the queue)
            -- receive (gets a message from the queue)
 		
        -t string body (argument for send - will send the string in this parameter in the message body
        -f file body (argument for send - will send the contents of the file parameter in the message body

    2.  ./sqsutils -a depth -s sqs.us-east-1.amazonaws.com -i 478989820108 DEV_TIBCO_STORES_QUEUE   
        // Talk to a real amazon queue

    3. ./sqsutils test_queue_1   
        // Gets the queue depth from local elastic mq queue test_queue_1
        
    4.  ./sqsutils -h
        //Prints basic help on the console

