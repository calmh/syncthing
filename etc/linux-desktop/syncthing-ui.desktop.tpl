[Desktop Entry]
Name=Syncthing Web UI
{{- range translate "Syncthing Web UI" }}
Name[{{.Lang}}]={{.Phrase}}
{{- end }}
GenericName=File synchronization UI
{{- range translate "File synchronization UI" }}
GenericName[{{.Lang}}]={{.Phrase}}
{{- end }}
Comment=Opens Syncthing's Web UI in the default browser (Syncthing must already be started).
{{- range translate "Opens Syncthing's Web UI in the default browser (Syncthing must already be started)." }}
Comment[{{.Lang}}]={{.Phrase}}
{{- end }}
Exec=/usr/bin/syncthing -browser-only
Icon=syncthing
Terminal=false
Type=Application
Keywords=synchronization;interface;
Categories=Network;FileTransfer;P2P
