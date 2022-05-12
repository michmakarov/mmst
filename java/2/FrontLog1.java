//220509 13:41 See ../HttpServer.java
//_______13:53 Let's this (FrontLog1) will be using by itself despite that I cannot formulate clearly what it means.
//220510 12:54 1. I indeed cannot understand the last phrase. Lo, I have remembered! This class would be checked for fitting by the owner through annotations.
//1. after https://techblogstation.com/java/java-write-to-file-line-by-line/
//220511 13:57 I have intended to make the Write(String msg) method static.
//There is suspicious that I have not enough knowledge to do so.
import java.io.File;
import java.io.FileOutputStream;
//import java.time;
import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;
import java.nio.charset.Charset;

public class FrontLog1 {
	//static private boolean isReady = false;
	private File file;	
	private FileOutputStream fileOutStream;
	private byte[] b;
	
        public FrontLog1(String fileName) throws Throwable {
			file = new File(fileName);
			file.createNewFile();
			fileOutStream = new FileOutputStream(file);
        }
        
    public void Write(String msg) throws Throwable {
 		try {
			msg = LocalDateTime.now() + " -> "+msg+"\n";
			b =  msg.getBytes("UTF-8");
			//System.out.printf("FrontLog1.Write: before writing; msg=%s\n", msg);
 			fileOutStream.write(b);
			fileOutStream.flush();
		}
		catch (Exception e) {
			e.printStackTrace();
		} finally {
			//try {
			//	if (fileOutStream != null) {
			//		fileOutStream.close();
			//	}
			//} catch (Exception e) {
			//	e.printStackTrace();
			//}
		}
	}
}
