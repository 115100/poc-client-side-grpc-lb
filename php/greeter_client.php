<?php
/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// php:generate protoc --proto_path=./../protos   --php_out=./   --grpc_out=./ --plugin=protoc-gen-grpc=./../../bins/opt/grpc_php_plugin ./../protos/helloworld.proto

require dirname(__FILE__).'/vendor/autoload.php';

@include_once dirname(__FILE__).'/Greeterpb/GreeterClient.php';
@include_once dirname(__FILE__).'/Greeterpb/GreetReply.php';
@include_once dirname(__FILE__).'/Greeterpb/GreetRequest.php';
@include_once dirname(__FILE__).'/GPBMetadata/Greeter.php';

function main()
{
    $client = new Greeterpb\GreeterClient(
        # headless service exposed here; always in fmt dns:///<svc>.<namespace>.svc.cluster.local:<port>
        'dns:///greeter-server.default.svc.cluster.local:8080', [
        # round_robin is a built-in balancer
        'grpc.lb_policy_name' => 'round_robin',
        'credentials' => Grpc\ChannelCredentials::createInsecure(),
    ]);

    while (1) {
        $request = new Greeterpb\GreetRequest();
        $request->setName('ping');
        list($reply, $status) = $client->Greet($request)->wait();
        if ($status->code != 0) {
            echo "client: failed to get response: ".$status->details."\n";
            return;
        }

        $message = $reply->getMessage();
        echo "client: got response from server: ".$message."\n";
        sleep(1);
    }
}

main();
