//220511 16:21 See ./HttpServer.java
import java.io.File;
import java.io.FileInputStream;

public class FileReader {
    public byte[] read(String fileName) throws Throwable {
		File file = new File(fileName);
		FileInputStream fileInStream = new FileInputStream(file);
		byte[] b = new byte[Globals.MaxFileSize];
 		try {
			//System.out.printf("FrontLog1.Write: before writing; msg=%s\n", msg);
 			fileOutStream.read(b);
			fileOutStream.flush();
		}
		catch (Throwable t) {
			t.printStackTrace();
		} finally {
		}
	}
}
