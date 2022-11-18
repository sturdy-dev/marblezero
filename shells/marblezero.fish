function marblezero_preexec --on-event fish_preexec
  marblezero --import-single "$argv"
end