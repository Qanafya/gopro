# Midterm 2:
Added new features: comment items, rating items, filtering data.<br />
A user can comment and rate any items on the sale. <br />
Also they can search item by filtering their price and rating.<br />
Explanation video: https://www.tiny.cc/qanafyamidterm <br />
If you have questions you can write in telegram @qanafya.<br /><hr><br />



Team member:
Yerali Ussen 200103214 (alone hero)<br />
 <br /><br />

Second progress: created database and register function

Third progress: imported database from open source codes, created login function. 

Fourth progress: added mysql database.

Fifth progress: saving a data without mysql, simple registration and authorization function.

Midterm 1: i couldn't push my midterm 1 to this repo, so i created new repository: https://github.com/Qanafya/gomid

Sixth progress: added products page. This page displays the products that are for sale.

Seventh and eighth progress: Added new product details page. And tested a insertions of values from database. 

Nineth progress: users can comment to the products and will record to db

Tenth progress: products page connected with database. It will show products from the db.




How to Install OpenSSH on Windows Server 2019 or 2022
1. Using Windows PowerShell
Open PowerShell as an Administrator.

2. 
Get-WindowsCapability -Online | Where-Object Name -like 'OpenSSH*'

3. Install OpenSSH Client:
Add-WindowsCapability -Online -Name OpenSSH.Client~~~~0.0.1.0

4. Install OpenSSH Server:

Add-WindowsCapability -Online -Name OpenSSH.Server~~~~0.0.1.0


5. Configure Firewall
- Open the Windows start menu, locate and click Server Manager. In the Server Manager window, navigate to Tools, and select Windows Defender Firewall with Advanced Security from the drop-down list.
- Now, click Inbound Rules in the open Firewall window. Then, select New Rule from the right pane.
In the New Inbound Rule Wizard, select Port from the list of options, then click Next. Select ‘TCP’, then enter port 22 in the Specific local ports: section.
- Next, allow the connection, assign the rule to server profiles, and set a custom name for easy identification from the list of Firewall rules.

Click Finish to save the new firewall rule.

6. Start-Service sshd
 restart-Service sshd

7. Set-Service -Name sshd -StartupType 'Automatic'

8. To configure OpenSSH, use the following command to open the main configuration file in Notepad and make your preferred changes.
 start-process notepad C:\Programdata\ssh\sshd_config


9. Login to Windows Server using SSH
$ ssh -l Administrator SERVER-IP -p port
