//220506 10:21
//220509 13:10 I have had an intention to release in first turn a front log.
//But now I have been diving into studying of annotations.
//Let's stand the task as:
//1. Let's the log is a distinct class FrontLog
//2. Let's build a mechanism that check the FrontLog is fit for using.
//220510 07:11 Why the second request is? Let's the problem name is the problem220510
//220511 12:57 My investigating of the problem220510 has been complicated by the strange behavior of a call of FrontLog1.Write("Start") in very begin of the main method.
//_______14:16 Now the FrontLog1.Write is non static
import java.net.ServerSocket;
import java.net.Socket;


public class HttpServer {
	
	private static FrontLog1 log;


    public static void main(String[] args) throws Throwable {
		log = new FrontLog1("FrontLog3.log");
		log.Write("Start");
		int clientCount=0;
        ServerSocket ss = new ServerSocket(8080);
        while (true) {
            System.out.printf("Waiting for socket(count=%d)\n", clientCount);
            Socket s = ss.accept();
            clientCount++;
            System.out.printf("________there is it (count=%d):RA=%s;KA=%s \n", clientCount, s.getInetAddress(), s.getKeepAlive());
            new Thread(new SocketProcessor(s, clientCount, log)).start();
        }
    }

}
