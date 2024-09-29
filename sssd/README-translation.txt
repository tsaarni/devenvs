

cat > mylang.po <<EOF
msgid ""
msgstr ""
"Language: mylang\n"
"Content-Type: text/plain; charset=UTF-8\n"
"Content-Transfer-Encoding: 8bit\n"

msgid "Your password has expired. You have %1$d grace login(s) remaining."
msgstr "Your password has expired. You have %1$d grace login(s) remaining. Your access will be revoked."
EOF

msgfmt mylang.po -o mylang.mo


sudo mkdir -p /proc/$(pidof sssd)/root/usr/share/locale/mylang/LC_MESSAGES/
sudo cp mylang.mo /proc/$(pidof sssd)/root/usr/share/locale/mylang/LC_MESSAGES/sssd.mo

sudo bash -c "echo 'LANGUAGE=mylang' >> /proc/$(pidof sssd)/root/etc/security/pam_env.conf"



sshpass -p expired ssh expired@localhost -p 2222 -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no "echo Hello world!"
