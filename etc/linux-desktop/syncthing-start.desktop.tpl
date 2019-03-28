[Desktop Entry]
Name=Start Syncthing
{{- range translate "Start Syncthing" }}
Name[{{.Lang}}]={{.Phrase}}
{{- end }}
GenericName=File synchronization
{{- range translate "File synchronization" }}
GenericName[{{.Lang}}]={{.Phrase}}
{{- end }}
Comment=Starts the main syncthing process in the background.
{{- range translate "Starts the main syncthing process in the background." }}
Comment[{{.Lang}}]={{.Phrase}}
{{- end }}
Exec=/usr/bin/syncthing -no-browser
Icon=syncthing
Terminal=false
Type=Application
Keywords=synchronization;daemon;
Categories=Network;FileTransfer;P2P
