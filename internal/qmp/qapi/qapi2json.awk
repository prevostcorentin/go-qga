BEGIN {
    print "["
    buffer = ""
    first = 1
    inblock = 0
}

/^[[:space:]]*#/ || /^[[:space:]]*$/ {
    next
}

/^\{/ {
    buffer = $0
    inblock = 1
    next
}

inblock {
    buffer = buffer " " $0
    if (/\}$/) {
        # Fin de bloc
        if (!first) {
            printf(",\n")
        }
        gsub(/'/, "\"", buffer)   # optionnel : remplacer ' par " pour JSON valide
        printf("%s", buffer)
        first = 0
        buffer = ""
        inblock = 0
    }
    next
}
END {
    print "\n]"
}
