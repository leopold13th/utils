package my.leo;

import java.io.File;
import java.io.IOException;
import java.util.ArrayList;
import java.util.Collection;
import java.util.Collections;

class FileHits {
    public String filename;
    public int hits;

    public FileHits(String filename, int hits) {
        this.filename = filename;
        this.hits = hits;
    }
}

public class Main
{
    static String path;
    static ArrayList<String> files = new ArrayList<>();
    static ArrayList<String> searchWords = new ArrayList<>();
    static ArrayList<FileHits> hitFiles = new ArrayList<>();

    public static void main( String[] args ) throws IOException {
        if (args.length > 1) {
            path = args[0];
            for (int i=1; i<args.length; i++) {
                searchWords.add(args[i].toLowerCase());
            }
        }
        tree(path);

        for (String fileName : files) {
            int hits = 0;
            for (String word : searchWords) {
                if (fileName.toLowerCase().contains(word)) hits++;
            }
            if (hits > 0) hitFiles.add(new FileHits(fileName, hits));
        }
        Collections.sort(hitFiles, (x,y) -> y.hits - x.hits);
        hitFiles.stream().forEach(x -> System.out.println(x.hits + " " + x.filename.substring(path.length())));
    }

    static void tree(String path) {
        File dir = new File(path);
        if (dir.isDirectory()) {
            for (File file : dir.listFiles()) {
                if (file.isDirectory()) {
                    tree(file.getAbsolutePath());
                } else {
                    files.add(file.getAbsolutePath());
                }
            }
        }
    }
}
