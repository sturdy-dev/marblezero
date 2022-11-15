function hackagotchi_preexec --on-event fish_preexec
  hackagotchi --import-single "$argv"
end