<?php
header("Content-Type: text/event-stream");
header("Cache-Control: no-cache");
header("Connection: keep-alive");
header("Access-Control-Allow-Origin: *");
header("Access-Control-Allow-Headers: Content-Type");
header("Access-Control-Allow-Methods: GET, OPTIONS");

if ($_SERVER['REQUEST_METHOD'] === 'OPTIONS') {
    http_response_code(200);
    exit;
}

$loremIpsum = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.";
$words = explode(" ", $loremIpsum);

$message = isset($_GET['message']) ? $_GET['message'] : '';

if (empty($message)) {
    echo "event: chatError\n";
    echo 'data: {"error": "Message query parameter is required"}\n\n';
    ob_flush();
    flush();
    return;
}

$sendImage = stripos($message, "image") !== false;

foreach ($words as $index => $word) {
    echo "event: message\n\n";
    $data = json_encode(['message' => $word]);
    echo "data: {$data}\n\n";
    ob_flush();
    flush();
    sleep(1);

    if ($index === count($words) - 1 && $sendImage) {
        $imgMessage = json_encode(['message' => 'Image URL: https://source.unsplash.com/random']);
        echo "data: {$imgMessage}\n\n";
        ob_flush();
        flush();
    }
}

echo "event: close\n";
echo "data: Done\n\n";
ob_flush();
flush();
