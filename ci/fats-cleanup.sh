#!/bin/bash

set -o nounset

fats_dir=`dirname "${BASH_SOURCE[0]}"`/fats

# attempt to cleanup fats
if [ -d "$fats_dir" ]; then
  echo "Uninstall riff system"
  helm delete riff --purge || true
  kubectl delete namespace $NAMESPACE || true

  source $fats_dir/cleanup.sh
fi
