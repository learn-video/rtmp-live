alias re := reload-edge

# reload NGINX
reload-edge:
    docker compose kill edge -s SIGHUP

# ffmpeg stream
stream:
    ffmpeg -re -f lavfi -i "smptehdbars=rate=30:size=1920x1080" \
        -f lavfi -i "sine=frequency=1000:sample_rate=48000" \
        -vf drawtext="text='RTMP Live %{localtime\:%X}':rate=30:x=(w-tw)/2:y=(h-lh)/2:fontsize=48:fontcolor=white:box=1:boxcolor=black" \
        -f flv -c:v h264 -profile:v baseline -pix_fmt yuv420p -preset ultrafast -tune zerolatency -crf 28 -g 60 -c:a aac \
        "rtmp://localhost:1935/stream/golive"
