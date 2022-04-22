package hello;

import common.print.Printer;
import hello.util.HelloUtil;

public class Hello {
  public static void main(String[] args) {
    Printer printer = new Printer();
    printer.print(HelloUtil.hello);
  }
}

