#!/bin/bash

kpxc_path="/path/to/database.kdbx"
function kpe() {
	temp_file=$(mktemp)
	exec 3> "$temp_file"
	kpxc-exp "-d=$kpxc_path" "$@" 3>&3
	exec 3>&-
	result=$(<"$temp_file")
	rm "$temp_file"

	IFS_old="$IFS"
	IFS='|&;'
	for secret in $result
	do
	    eval "export $secret > /dev/null" 
	done
	IFS="$IFS_old"
}