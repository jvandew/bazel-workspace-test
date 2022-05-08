package hello;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import common.print.Printer;
import hello.util.HelloUtil;

public class Hello {
  public static void main(String[] args) throws JsonProcessingException {
    ObjectMapper mapper = new ObjectMapper();
    Hello hello = mapper.readValue("{}", Hello.class);

    Printer printer = new Printer();
    printer.print(HelloUtil.hello);
  }
}
