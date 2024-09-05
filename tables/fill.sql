INSERT INTO State (name)
VALUES ("init") ,  ("downloading"), ("done"), ("failed"), ("paused"), ("canceled");

INSERT INTO DownloadType (name)
VALUES ("http"), ("ftp"), ("sftp"), ("sftp+http"), ("sftp+ftp");

INSERT INTO FileType (name)
VALUES ("undefined"), ("text"), ("binary"), ("image"), ("audio"), ("video");

INSERT INTO SETTINGS 
VALUES ("./tmp/working", "./tmp/output");
