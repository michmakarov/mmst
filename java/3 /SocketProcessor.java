import java.net.ServerSocket;
import java.net.Socket;
import java.io.InputStream;
import java.io.OutputStream;
import java.io.InputStreamReader;
import java.io.BufferedReader;


    class SocketProcessor implements Runnable {

        private Socket s;
        private InputStream is;
        private OutputStream os;
        private int socketNum;
		private FrontLog1 log;

        public SocketProcessor(Socket s, int num, FrontLog1 log) throws Throwable {
            this.s = s;
            this.is = s.getInputStream();
            this.os = s.getOutputStream();
            this.socketNum = num;
            this.log = log;
        }

        public void run() {
            try {
				WriteToFL(s, socketNum);
                //readInputHeaders();
                writeResponse("<html><body><h1>The main page of the mmst будет тута!</h1></body></html>");
            } catch (Throwable t) {
                /*do nothing*/
            } finally {
                try {
                    s.close();
                } catch (Throwable t) {
                    /*do nothing*/
                }
            }
            System.err.println("Client processing finished");
        }
        
        private void WriteToFL(Socket s, int socketNum) {
			try {
			log.Write("(" + socketNum + ")" + s.toString());
			} catch (Throwable t){
				try{
					log.Write("(" + socketNum + ")Err:" + t.getLocalizedMessage());
				} catch(Throwable tt) {}
			}
		}

        private void writeResponse(String s) throws Throwable {
            String response = "HTTP/1.1 200 OK\r\n" +
                    "Server: YarServer/2009-09-09\r\n" +
                    "Content-Type: text/html; charset=utf-8\r\n" +
                    "Content-Length: " + s.length() + "\r\n" +
                    "Connection: close\r\n\r\n";
            String result = response + s;
            os.write(result.getBytes());
            os.flush();
        }

        private void readInputHeaders() throws Throwable {
			int count=0;
            BufferedReader br = new BufferedReader(new InputStreamReader(is));
            while(true) {
				count++;
                String s = br.readLine();
                System.out.printf("%d)%s\n", count, s);
                if(s == null || s.trim().length() == 0) {
                    break;
                }
            }
        }
    }

