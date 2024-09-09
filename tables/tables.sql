CREATE TABLE IF NOT EXISTS State (
    ID_State INTEGER PRIMARY KEY AUTOINCREMENT, 
    Name TEXT
);

CREATE TABLE IF NOT EXISTS DownloadType (
    ID_Download_Type INTEGER PRIMARY KEY AUTOINCREMENT,
    Name TEXT DEFAULT "http"
);

CREATE TABLE IF NOT EXISTS FileType (
    ID_File_Type INTEGER PRIMARY KEY AUTOINCREMENT,
    Name TEXT DEFAULT "undefined"
);

CREATE TABLE IF NOT EXISTS Packet (
    ID_Packet INTEGER PRIMARY KEY AUTOINCREMENT,
    Start INTEGER,
    End INTEGER,
    ID_Packet_State INTEGER,
    ID_Download INTEGER, 
    FOREIGN KEY (ID_Packet_State) REFERENCES State(ID_State)
);

CREATE TABLE IF NOT EXISTS Download (
    ID_Download INTEGER PRIMARY KEY AUTOINCREMENT,
    ID_Download_Type INTEGER,
    ID_Download_State INTEGER,
    ID_File_Type INTEGER,
    Working_file_path TEXT,
    Output_file_path TEXT,
    Remote TEXT,
    FOREIGN KEY (ID_Download_Type) REFERENCES DownloadType(ID_Download_Type),
    FOREIGN KEY (ID_Download_State) REFERENCES State(ID_State),
    FOREIGN KEY (ID_File_Type) REFERENCES FileType(ID_File_Type)
);


CREATE TABLE IF NOT EXISTS settings (
    Working_dir TEXT ,
    Output_dir TEXT , 
    PacketSize INTEGER 
);


